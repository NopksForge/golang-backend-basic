package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
	"training/app/product"
	"training/app/user"
	"training/logger"
	"training/redis"

	"training/app"
	"training/config"
	"training/database"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"

	_ "embed"
	_ "time/tzdata"
)

const (
	gracefulShutdownDuration = 10 * time.Second
	serverReadHeaderTimeout  = 5 * time.Second
	serverReadTimeout        = 5 * time.Second
	serverWriteTimeout       = 10 * time.Second // request hangup after this durations
	handlerTimeout           = serverWriteTimeout - (time.Millisecond * 100)
)

// go build -ldflags "-X main.commit=123456"
var commit string

//go:embed VERSION
var version string

func init() {
	if os.Getenv("GOMAXPROCS") != "" {
		runtime.GOMAXPROCS(0) // GOMAXPROCS
	} else {
		runtime.GOMAXPROCS(1) // 0 - 999m
	}
	if os.Getenv("GOMEMLIMIT") != "" {
		debug.SetMemoryLimit(-1) // GOMEMLIMIT
	}
}

func main() {
	cfg := config.C(config.Env)
	_ = logger.New(logger.GCPKeyReplacer)

	r, stop := router(cfg)
	defer stop()

	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           r,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		MaxHeaderBytes:    1 << 20,
	}

	go gracefully(srv)

	slog.Info("run at :" + cfg.Server.Port)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("HTTP server ListenAndServe: " + err.Error())
		return
	}

	slog.Info("bye")
}

func router(cfg config.Config) (*gin.Engine, func()) {
	r := gin.New()
	r.Use(gin.Recovery())

	if config.IsDevEnv() {
		r.Use(gin.Logger())
	}

	// health check handler
	{
		r.GET("/liveness", liveness())
		r.GET("/metrics", metrics())
		r.GET("/readiness", readiness())
	}

	r.Use(
		securityHeaders,
		accessControl,
		app.RefIDMiddleware(cfg.Header.RefIDHeaderKey),
		app.AutoLoggingMiddleware,
		handlerTimeoutMiddleware,
	)

	apiV1Router := r.Group("/api/v1")

	// kafka
	kafkaServer := []string{"localhost:9092"}
	kafkaProducer, err := sarama.NewSyncProducer(kafkaServer, nil)
	if err != nil {
		panic(err)
	}
	kafkaConsumer, err := sarama.NewConsumer(kafkaServer, nil)
	if err != nil {
		panic(err)
	}

	db := database.NewPostgresDB(cfg.Database.PostgresURL)
	redisClient := redis.New(cfg.Database.RedisURL, "")

	// app/product
	productRouter := apiV1Router.Group("/product")
	{
		productRepository := product.NewRepository(db, redisClient)
		productService := product.NewService(productRepository)
		productHandler := product.NewHandler(productService)
		productHandler.InitEndpoints(productRouter)
	}

	// app/user
	userRouter := apiV1Router.Group("/user")
	{
		userRepository := user.NewRepository(user.NewRepositoryCfg{
			Db:            db,
			Redis:         redisClient,
			KafkaProducer: kafkaProducer,
		})
		userService := user.NewService(userRepository, kafkaConsumer)
		userHandler := user.NewHandler(userService)
		userHandler.InitEndpoints(userRouter)

		if err := userService.ConsumeUserCreation(context.Background()); err != nil {
			panic(err)
		}
	}
	return r, func() {
		db.Close()
		redisClient.Close()
		kafkaProducer.Close()
		kafkaConsumer.Close()
	}
}

func handlerTimeoutMiddleware(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerTimeout)
	defer cancel()
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func securityHeaders(c *gin.Context) {
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Header("Referrer-Policy", "strict-origin")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
	c.Next()
}

var headers = []string{
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"X-CSRF-Token",
	"Authorization",
	"accept",
	"origin",
	"Cache-Control",
	"X-Requested-With",
	os.Getenv("REF_ID_HEADER_KEY"),
}

func accessControl(c *gin.Context) {
	cfg := config.C(config.Env)
	c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.AccessControl.AllowOrigin)
	c.Writer.Header().Set("Access-Control-Request-Method", "POST, GET, PUT, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

func gracefully(srv *http.Server) {
	{
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()
		<-ctx.Done()
	}

	d := time.Duration(gracefulShutdownDuration)
	slog.Info(fmt.Sprintf("shutting down in %d ...\n", d))
	// We received an interrupt signal, shut down.
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		slog.Info("HTTP server Shutdown: " + err.Error())
	}
}

func liveness() func(c *gin.Context) {
	h, err := os.Hostname()
	if err != nil {
		h = fmt.Sprintf("unknown host err: %s", err.Error())
	}
	return func(c *gin.Context) {
		// the liveness probe is only this API itself probe
		// others service healthy not responsibility of this API
		// however, if you need it please follow these steps yourself
		// - check db connection, redis connection, etc
		// - implement help check your service dependencies
		// - implement help check for Postgres, MongoDB, Redis, etc
		//   e.g. MongoDB database.IsMongoReady()
		//   e.g. Redis database.IsRedisReady()
		//   e.g. Kafka database.IsKafkaReady()

		// e.g. check if Postgres is ready
		// if !database.IsPostgresReady() {
		// 	c.Status(http.StatusInternalServerError)
		// 	return
		// }

		c.JSON(http.StatusOK, gin.H{
			"hostname": h,
			"version":  strings.ReplaceAll(version, "\n", ""),
			"commit":   commit,
		})
	}
}

func readiness() gin.HandlerFunc {
	// when the server go to ListenAndServe it means everything is ready to serve
	// no need and checking such as db checking again

	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func metrics() func(c *gin.Context) {
	return func(c *gin.Context) {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		c.JSON(http.StatusOK, gin.H{
			"memory": gin.H{
				"alloc":        toMB(mem.Alloc),
				"totalAlloc":   toMB(mem.TotalAlloc),
				"sysAlloc":     toMB(mem.Sys),
				"heapInuse":    toMB(mem.HeapInuse),
				"heapIdle":     toMB(mem.HeapIdle),
				"heapReleased": toMB(mem.HeapReleased),
				"stackInuse":   toMB(mem.StackInuse),
				"stackSys":     toMB(mem.StackSys),
			},
		})
	}
}

type Size uint64

const (
	Byte Size = 1 << (10 * iota)
	KB
	MB
)

func megabytes(b uint64) float64 {
	return float64(b) / float64(MB)
}

func toMB(b uint64) string {
	return fmt.Sprintf("%.2f MB", megabytes(b))
}

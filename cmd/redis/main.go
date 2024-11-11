package main

import (
	"context"
	"fmt"
	"time"
	"training/redis"
)

func main() {
	addr := "localhost:6379"
	redisClient := redis.New(addr, "")
	defer redisClient.Close()

	ctx := context.Background()
	if err := redisClient.Set(ctx, "name", "ArmLnwZa", 1*time.Minute).Err(); err != nil {
		panic(err)
	}

	val, err := redisClient.Get(ctx, "name").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(val)
}

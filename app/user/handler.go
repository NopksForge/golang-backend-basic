package user

import (
	"errors"
	"net/http"
	"training/app"
	"training/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitEndpoints(r gin.IRoutes) {
	r.POST("/", h.Save)
	r.PUT("/:userId", h.Update)
	r.DELETE("/:userId", h.Delete)
	r.GET("/:userId", h.Get)
}

func (h *Handler) Get(c *gin.Context) {
	var req GetUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, app.Response{Message: err.Error()})
		return
	}

	if err := validator.Validate(c, req); err != nil {
		c.JSON(http.StatusBadRequest, app.Response{Message: err.Error()})
		return
	}

	user, err := h.service.GetById(c, uuid.MustParse(req.UserId))
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			c.JSON(http.StatusNotFound, app.Response{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, app.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, GetUserResponse{
		UserId:    user.UserId.String(),
		UserEmail: user.UserEmail,
		UserName:  user.UserName,
	})
}

func (h *Handler) Save(c *gin.Context) {
	var req SaveUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, app.Response{Message: err.Error()})
		return
	}

	if err := validator.Validate(c, req); err != nil {
		c.JSON(http.StatusBadRequest, app.Response{Message: err.Error()})
		return
	}

	payload, err := h.service.Save(c, SaveUserPayload{
		UserEmail: req.UserEmail,
		UserName:  req.UserName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, app.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, payload)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, app.Response{Message: err.Error()})
		return
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, app.Response{Message: err.Error()})
		return
	}

	if err := validator.Validate(c, req); err != nil {
		c.JSON(http.StatusBadRequest, app.Response{Message: err.Error()})
		return
	}

	err := h.service.Update(c, UpdateUserPayload{
		UserId:    uuid.MustParse(req.UserId),
		UserEmail: req.UserEmail,
		UserName:  req.UserName,
	})
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			c.JSON(http.StatusNotFound, app.Response{Message: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, app.Response{Message: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	var req DeleteUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, app.Response{Message: err.Error()})
		return
	}

	if err := validator.Validate(c, req); err != nil {
		c.JSON(http.StatusBadRequest, app.Response{Message: err.Error()})
		return
	}

	err := h.service.Delete(c, uuid.MustParse(req.UserId))
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			c.JSON(http.StatusNotFound, app.Response{Message: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, app.Response{Message: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

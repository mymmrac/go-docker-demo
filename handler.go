package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	log *zap.SugaredLogger
	e   *echo.Echo
}

func NewHandler(log *zap.SugaredLogger, e *echo.Echo) *Handler {
	return &Handler{
		log: log,
		e:   e,
	}
}

func (h *Handler) RegisterRoutes() {
	h.e.GET("/", h.index)
	h.e.GET("/health", h.health)
}

func (h *Handler) index(ctx echo.Context) error {
	n := rand.Int()

	if n%2 == 0 {
		h.log.Info("N is even")
	} else {
		h.log.Info("N is odd")
	}

	return ctx.JSON(http.StatusOK, "Hello, World!")
}

type health struct {
	Running bool      `json:"running"`
	Time    time.Time `json:"time"`
}

func (h *Handler) health(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, health{
		Running: true,
		Time:    time.Now().UTC(),
	})
}

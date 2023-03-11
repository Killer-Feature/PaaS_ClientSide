package handlers

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"KillerFeature/ClientSide/internal"
)

type Handler struct {
	logger *zap.Logger
	u      internal.Usecase
}

func NewHandler(logger *zap.Logger, u internal.Usecase) *Handler {
	return &Handler{logger: logger, u: u}
}

func (h *Handler) Register(s *echo.Echo) {
	// Register http handlers
	s.GET("/hello", h.GetHello)
}

func (h *Handler) GetHello(c echo.Context) error {
	return c.HTML(http.StatusOK, "hello")
}

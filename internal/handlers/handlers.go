package handlers

import (
	"embed"
	"io/fs"
	"log"
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

//go:embed dist
var ui embed.FS

func (h *Handler) Register(s *echo.Echo) {
	// Register http handlers
	s.GET("/hello", h.GetHello)

	fsys, err := fs.Sub(ui, "dist")
	if err != nil {
		log.Fatal(err)
	}

	s.GET("/*", echo.WrapHandler(http.FileServer(http.FS(fsys))))
}

func (h *Handler) GetHello(c echo.Context) error {
	h.logger.Info("request received", zap.String("host", c.Request().RemoteAddr))
	return c.HTML(http.StatusOK, "hello")
}

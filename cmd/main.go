package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/netip"

	"github.com/labstack/echo/v4/middleware"

	"github.com/Killer-Feature/PaaS_ClientSide/pkg/helm"
	k8s_installer "github.com/Killer-Feature/PaaS_ClientSide/pkg/k8s-installer"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/logger/zaplogger"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/Killer-Feature/PaaS_ClientSide/internal/handlers"
	"github.com/Killer-Feature/PaaS_ClientSide/internal/repository"
	"github.com/Killer-Feature/PaaS_ClientSide/internal/service"
)

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, err := zaplogger.NewZapLogger(&config)
	servLogger := servlog.NewServLogger(prLogger)
	if err != nil {
		log.Fatal("zap logger build error")
	}
	logger := prLogger.Desugar()
	defer func(prLogger *zap.Logger) {
		err = prLogger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(logger)

	server := echo.New()

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	r, err := repository.Create(logger)
	if err != nil {
		logger.Fatal("database creating error", zap.Error(err))
	}
	tm := taskmanager.NewTaskManager[netip.AddrPort](ctx, servLogger)
	k8sinstaller := k8s_installer.NewInstaller(logger, r)
	hi, err := helm.NewHelmInstaller("default", "https://charts.bitnami.com/bitnami", "bitnami", logger)
	if err != nil {
		logger.Fatal("helm installer creating error", zap.Error(err))
	}
	u := service.NewService(r, logger, tm, k8sinstaller, hi)
	h := handlers.NewHandler(logger, u)
	h.Register(server)

	// metrics := monitoring.RegisterMonitoring(server)

	//m := middleware.NewMiddleware(p.Logger, metrics)
	//m.Register(server)
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXCSRFToken},
		AllowMethods:     []string{echo.GET, echo.POST, echo.OPTIONS, echo.PUT},
		ExposeHeaders:    []string{echo.HeaderXCSRFToken},
		MaxAge:           86400,
	}))

	g.Go(func() error {
		return server.Start(":80")
	})

	if err := g.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error shutdown with error", zap.Error(err))
			ctx := context.Background()
			//u.CloseService()
			_ = server.Shutdown(ctx)
		}
	}

}

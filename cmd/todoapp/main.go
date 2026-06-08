package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Astek27/todoapp/internal/core/logger"
	core_http_middleware "github.com/Astek27/todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Astek27/todoapp/internal/core/transport/http/server"
	users_transport_http "github.com/Astek27/todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("start application")

	usersTransport := users_transport_http.NewUsersHTTPHandler(nil)
	routes := usersTransport.Routes()

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(routes...)

	server := core_http_server.NewServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	server.RegisterRouters(apiVersionRouter)

	if err := server.Run(ctx); err != nil {
		logger.Error("HTTP server error", zap.Error(err))
	}
}
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/Astek27/todoapp/internal/core/logger"
	core_pgx_pool "github.com/Astek27/todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Astek27/todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Astek27/todoapp/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/Astek27/todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Astek27/todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/Astek27/todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Astek27/todoapp/internal/features/users/repository/postgres"
	users_service "github.com/Astek27/todoapp/internal/features/users/service"
	users_transport_http "github.com/Astek27/todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

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
	
	logger.Debug("application timezone", zap.Any("timezone", timeZone))

	logger.Debug("Initiialization postgres connection pool")

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("create connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("start feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransport := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("start feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransport := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("create server")
	
	server := core_http_server.NewServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransport.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransport.Routes()...)
	
	// apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(
	// 	core_http_server.ApiVersion2,
	// 	core_http_middleware.Dummy("api v2 middleware"),
	// )
	// apiVersionRouterV2.RegisterRoutes(usersTransport.Routes()...)

	server.RegisterRouters(apiVersionRouterV1)
	// server.RegisterRouters(apiVersionRouterV2)

	if err := server.Run(ctx); err != nil {
		logger.Error("HTTP server error", zap.Error(err))
	}
}
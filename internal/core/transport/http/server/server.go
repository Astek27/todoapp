package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/Astek27/todoapp/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger
}

func NewServer(config Config, log *core_logger.Logger) *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(),
		config: config,
		log: log,
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr: h.config.Addr,
		Handler: h.mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		h.log.Warn("Start HTTP server", zap.String("Addr", h.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()
	
	select {
	case <- ctx.Done():
		h.log.Warn("HTTP server shutdown...")

		shutDownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutDownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown http server: %w", err)
		}

	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve http: %w", err)
		}

	h.log.Warn("HTTP server stoped")
	}

	return nil
}
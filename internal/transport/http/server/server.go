package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	"go.uber.org/zap"
)

type HttpServer struct {
	mux        *http.ServeMux
	config     Config
	log        *core_logger.Logger
	middleware []http_middleware.Middleware
}

func NewHttpServer(config Config, log *core_logger.Logger, middleware ...http_middleware.Middleware) *HttpServer {
	return &HttpServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (h *HttpServer) Run(ctx context.Context) error {
	mux := http_middleware.ChainMiddleware(h.mux, h.middleware...)
	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)

		h.log.Warn("starting http server", zap.String("addr", h.config.Addr))
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve http: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutting down http server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), h.config.ShutdownTimeout)
		defer cancel()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			server.Close()
			return fmt.Errorf("shutdown http: %w", err)
		}
		h.log.Warn("http server stopped")
	}

	return nil
}

func (h *HttpServer) RegisterRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)
		h.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

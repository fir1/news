package http

import (
	"context"
	"fmt"
	"github.com/fir1/news/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
	"strings"
	"sync"
)

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

//	@title			Swagger News API
//	@version		1.0
//	@description	This is an API documentation for News Feed backend
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url
//	@contact.email	kasimovfirdavs@gmail.com

// @host		localhost:8080/
// @BasePath	/
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func (s *Service) Run(stop chan struct{}) error {
	s.router = chi.NewRouter()
	s.router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:     []string{"*"}, //TODO: must be changed to allow only prod, dev hosts.
			AllowedMethods:     []string{"GET", "POST", "HEAD", "PATCH", "OPTIONS", "GET", "PUT"},
			AllowedHeaders:     []string{"*"},
			ExposedHeaders:     nil,
			AllowCredentials:   true,
			MaxAge:             300,
			OptionsPassthrough: false,
			Debug:              true,
		}),
		middleware.Logger,
	)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", stripProtocol(s.config.ServerHostName), s.config.LoadBalancerHostPort)

	if s.config.LoadBalancerHostPort == 443 {
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s:%d/swagger/doc.json", s.config.ServerHostName, s.config.LoadBalancerHostPort)),
		httpSwagger.DomID("swagger-ui")),
	)

	// Register all routes on http handler
	s.routes()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.router,
	}

	// channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				var errPanic error
				//nolint:errorlint //false positive
				s.logger.Error(fmt.Errorf("%+v", errPanic))
			}
			wg.Done()
		}()
		s.logger.Printf("REST API listening on port: %d for environment: %s", s.config.Port, s.config.Environment)
		serverErrors <- server.ListenAndServe()
	}()

	// blocking run and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting REST API http: %w", err)
	case <-stop:
		s.logger.Warn("http receive STOP signal")
		// asking listener to shutdown
		err := server.Shutdown(context.Background())
		if err != nil {
			return fmt.Errorf("graceful shutdown did not complete: %w", err)
		}
		s.logger.Info("http was shut down gracefully")
	}
	return nil
}

func stripProtocol(url string) string {
	if strings.HasPrefix(url, "http://") {
		return strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		return strings.TrimPrefix(url, "https://")
	}
	return url
}

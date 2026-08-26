package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"personae-fasti/configs"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"

	"github.com/rs/cors"
)

type APIServer struct {
	Server             *http.Server
	router             *http.ServeMux
	auth               *service.AuthService
	imageUploadTimeout time.Duration
}

type HandlerFunc[Body any] func(httputils.RequestData[Body]) httputils.Responder

//type AuthHandlerFunc[Body any] func(PlayerAuth, RequestData[Body]) Response

// type PlayerAuth struct {
// 	userID int64
// }

type ReadinessCheck func(context.Context) error

func ConfigServer(c *configs.Main, s *service.AuthService, readinessCheck ReadinessCheck) *APIServer {

	router := http.NewServeMux()
	configureHealthEndpoints(router, readinessCheck)

	allowedOrigins := c.App.AllowedOrigins
	if len(allowedOrigins) == 0 {
		if c.App.Debug || !strings.EqualFold(c.App.Environment, "production") {
			allowedOrigins = []string{"*"}
		} else {
			allowedOrigins = []string{"https://app.storyshard.ru"}
		}
	}
	allowCredentials := c.App.AllowCredentials && !containsWildcard(allowedOrigins)

	crs := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: allowCredentials,
	})

	imageUploadTimeout := durationOrDefault(c.App.ImageUploadTimeout, 120*time.Second)
	readTimeout := durationOrDefault(c.App.ReadTimeout, imageUploadTimeout)
	writeTimeout := durationOrDefault(c.App.WriteTimeout, imageUploadTimeout)
	if readTimeout < imageUploadTimeout {
		readTimeout = imageUploadTimeout
	}
	if writeTimeout < imageUploadTimeout {
		writeTimeout = imageUploadTimeout
	}

	api := &APIServer{
		Server: &http.Server{
			Addr:              c.App.Port,
			Handler:           crs.Handler(router),
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       120 * time.Second,
		},
		router:             router,
		auth:               s,
		imageUploadTimeout: imageUploadTimeout,
	}

	return api
}

func configureHealthEndpoints(router *http.ServeMux, readinessCheck ReadinessCheck) {
	router.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	router.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if readinessCheck != nil {
			if err := readinessCheck(ctx); err != nil {
				http.Error(w, "not ready", http.StatusServiceUnavailable)
				return
			}
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})
}

func containsWildcard(origins []string) bool {
	for _, origin := range origins {
		if origin == "*" {
			return true
		}
	}
	return false
}

func durationOrDefault(seconds int, fallback time.Duration) time.Duration {
	if seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}

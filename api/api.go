package api

import (
	"net/http"
	"time"

	"personae-fasti/configs"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"

	"github.com/rs/cors"
)

type APIServer struct {
	Server *http.Server
	router *http.ServeMux
	auth   *service.AuthService
	// storage    *data.Storage
	// fileServer *configs.FileServer
	// auth       *configs.AuthConfig
}

type HandlerFunc[Body any] func(httputils.RequestData[Body]) httputils.Responder

//type AuthHandlerFunc[Body any] func(PlayerAuth, RequestData[Body]) Response

// type PlayerAuth struct {
// 	userID int64
// }

func ConfigServer(c *configs.Main, s *service.AuthService, l *service.LogService) *APIServer {

	router := http.NewServeMux()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	api := &APIServer{
		Server: &http.Server{
			Addr:         c.App.Port,
			Handler:      crs.Handler(router),
			ReadTimeout:  time.Duration(c.App.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(c.App.WriteTimeout) * time.Second,
		},
		router: router,
		auth:   s,
	}

	return api
}

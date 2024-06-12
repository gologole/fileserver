package server

import (
	"context"
	"main.go/config"
	"net/http"
)

type Server struct {
	server *http.Server
}

func (s *Server) RunServer(handler http.Handler) error {
	serverConfig := config.Conf.Server.HTTP

	s.server = &http.Server{
		Addr:         ":" + serverConfig.Port,
		Handler:      handler,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
	}
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

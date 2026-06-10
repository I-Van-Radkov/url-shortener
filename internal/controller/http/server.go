package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

func NewServer(port int, readTimeout, writeTimeout time.Duration, router http.Handler) *Server {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      router,
	}

	return &Server{
		srv: srv,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

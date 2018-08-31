package server

import (
	"net/http"
	"piaqua/pkg/controller"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	srv *http.Server
}

func (s *Server) Start(c *controller.Controller) error {
	router := httprouter.New()
	router.GET("/state", state(c))

	s.srv = &http.Server{Addr: ":8080", Handler: router}
	err := s.srv.ListenAndServe()
	return err
}

func (s *Server) Stop() {
	if s.srv != nil {
		s.srv.Close()
	}
}

package server

import (
	"log"
	"net"
	"net/http"
	"piaqua/pkg/controller"
	"time"
)

type HTTPServer struct {
	srv http.Server
}

func NewHTTPServer(c *controller.Controller) *HTTPServer {
	return &HTTPServer{
		srv: http.Server{
			Addr:           "0.0.0.0:80",
			Handler:        newHandler(c),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   7 * time.Second,
			IdleTimeout:    70 * time.Second,
			MaxHeaderBytes: 1 << 16,
		}}
}

func (s *HTTPServer) ListenAndServe() error {
	ln, err := net.Listen("tcp4", s.srv.Addr)
	if err != nil {
		return err
	}
	log.Println("http: Server listen on", ln.Addr())
	return s.srv.Serve(ln)
}

func (s *HTTPServer) Close() {
	s.srv.Close()
}

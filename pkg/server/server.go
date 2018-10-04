package server

import (
	"fmt"
	"net"
	"net/http"
	"piaqua/pkg/config"
	"piaqua/pkg/controller"
	"time"
)

type HTTPServer struct {
	srv http.Server
}

func NewHTTPServer(configDir string, c *controller.Controller) (*HTTPServer, error) {
	var conf config.ServerConf
	err := conf.Read(configDir)
	if err != nil {
		return nil, err
	}
	return &HTTPServer{
		srv: http.Server{
			Addr:           conf.Address,
			Handler:        newHandler(&conf, c),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   7 * time.Second,
			IdleTimeout:    70 * time.Second,
			MaxHeaderBytes: 1 << 16,
		}}, nil
}

func (s *HTTPServer) ListenAndServe() error {
	ln, err := net.Listen("tcp4", s.srv.Addr)
	if err != nil {
		return err
	}
	fmt.Println("http: Server listen on", ln.Addr())
	return s.srv.Serve(ln)
}

func (s *HTTPServer) Close() {
	s.srv.Close()
}

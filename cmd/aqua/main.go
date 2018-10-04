package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"piaqua/pkg/controller"
	"piaqua/pkg/server"
	"piaqua/pkg/singleinstance"
	"syscall"
)

const defaultConfigDir = "configs"

func main() {
	configDir := flag.String("c", defaultConfigDir, "config directory")
	flag.Parse()

	if !singleinstance.Lock("piaqua") {
		panic("Another instance of a program is already running")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctrl, err := controller.NewController(*configDir)
	if err != nil {
		panic(err)
	}
	fmt.Println("Controller started")

	httpServer, err := server.NewHTTPServer(*configDir, ctrl)
	if err != nil {
		panic(err)
	}

	go func() {
		<-quit
		httpServer.Close()
	}()

	fmt.Println(httpServer.ListenAndServe())

	ctrl.Stop()
	fmt.Println("Controller stopped")
}

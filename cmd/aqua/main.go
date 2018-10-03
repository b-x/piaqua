package main

import (
	"log"
	"os"
	"os/signal"
	"piaqua/pkg/controller"
	"piaqua/pkg/server"
	"piaqua/pkg/singleinstance"
	"syscall"
)

const configDir = "."

func main() {
	if !singleinstance.Lock("piaqua") {
		log.Fatalln("Another instance of a program is already running")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctrl, err := controller.NewController(configDir)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Controller started")

	httpServer := server.NewHTTPServer(ctrl)

	go func() {
		<-quit
		httpServer.Close()
	}()

	log.Println(httpServer.ListenAndServe())

	ctrl.Stop()
	log.Println("Controller stopped")
}

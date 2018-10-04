package main

import (
	"flag"
	"log"
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
		log.Fatalln("Another instance of a program is already running")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctrl, err := controller.NewController(*configDir)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Controller started")

	httpServer, err := server.NewHTTPServer(*configDir, ctrl)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		<-quit
		httpServer.Close()
	}()

	log.Println(httpServer.ListenAndServe())

	ctrl.Stop()
	log.Println("Controller stopped")
}

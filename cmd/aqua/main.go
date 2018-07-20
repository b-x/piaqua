package main

import (
	"log"
	"os"
	"os/signal"
	"piaqua/pkg/controller"
	"piaqua/pkg/singleinstance"
	"syscall"
)

const configDir = "."

func main() {
	if !singleinstance.Lock("piaqua") {
		log.Fatalln("Another instance of a program is already running")
	}

	quit := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		close(done)
	}()

	log.Println("Controller is starting")
	c, err := controller.NewController(configDir)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Controller started")
	<-done
	log.Println("Controller is stopping")
	c.Stop()
	log.Println("Controller stopped")
}

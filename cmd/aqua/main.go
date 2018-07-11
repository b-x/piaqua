package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quit := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		close(done)
	}()

	log.Println("Controller is starting")
	<-done
	log.Println("Controller stopped")
}

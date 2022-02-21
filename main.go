package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

//
var build = "develop"

func main() {
	// how many CPU cores would be run in parallel
	g := runtime.GOMAXPROCS(0)
	log.Printf("starting service %q CPUS {%d}", build, g)
	defer log.Println("service ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("stopping servicea")
}

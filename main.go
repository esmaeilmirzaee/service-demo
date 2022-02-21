package main

import (
	"fmt"
	"go.uber.org/automaxprocs/maxprocs"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "go.uber.org/automaxprocs/maxprocs"
)

//
var build = "develop"

func main() {
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("maxprocs %w", err)
		os.Exit(1)
	}

	// how many CPU cores would be run in parallel
	g := runtime.GOMAXPROCS(0)

	log.Printf("starting service %q CPUS {%d}", build, g)
	defer log.Println("service ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("stopping servicea")
}

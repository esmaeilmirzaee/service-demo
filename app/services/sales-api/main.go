package main

import (
	"errors"
	"expvar"
	"fmt"
	"github.com/esmaeilmirzaee/service/app/services/sales-api/handlers"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
)

/*
TODO: Need to figure out timeouts for http service
*/

var build = "develop"

func main() {
	log, err := initLogger("SALES-API")
	if err != nil {
		fmt.Println("Error constructing logger", err)
		os.Exit(1)
	}

	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {
			log.Errorw("deferring logger", "ERROR", err)
			os.Exit(1)
		}
	}(log)

	// Perform the startup and shutdown sequence
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// ====================================================
	// GOMAXPROXS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotes
	undo, err := maxprocs.Set()
	defer undo()
	if err != nil {
		return fmt.Errorf("maxprcs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// ====================================================
	// configuration

	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:1337"`
			DebugHost       string        `conf:"default:0.0.0.0:1338"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
		}
	}{
		Version: conf.Version{
			SVN:  build,
			Desc: "copyright information here",
		},
	}

	const prefix = "SALES"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config %w", err)
	}

	// ==========================================================
	// App starting
	log.Infow("Starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// Display cmd information
	expvar.NewString("build").Set(build)

	// ==========================================================
	// Start Debug Service
	log.Infow("startup", "status", "debug router started", "host",
		cfg.Web.DebugHost)

	// The Debug function returns a mux to listen and serve on for
	// all the debug related endpoints. This include the standard
	// library endpoints.

	// construct the mux for the debug calls.
	debugMux := handlers.DebugStandardLibraryMux()

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding.
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host",
				cfg.Web.DebugHost, "Error", err)
		}
	}()

	// ==========================================================
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	return nil
}

func initLogger(service string) (*zap.SugaredLogger, error) {
	// Construct the application human-readable logger.
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": "SALES-API",
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}

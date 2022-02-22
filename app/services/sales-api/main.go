package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

//
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

	log,err:=config.Build()
	if err!= nil {
		return nil, err
	}


	return log.Sugar(), nil
}

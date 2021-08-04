package main

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	config, err := ParseConfigFromFlags()
	if err != nil {
		check.ExitError(err)
	}

	if config.Debug {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	}

	api := RestAPI{URL: config.API, Client: config.NewClient()}

	result, err := api.ExecuteCheck(config.Command, config.Arguments)
	if err != nil {
		check.ExitError(err)
	}

	_, _ = fmt.Fprintln(os.Stdout, result.String())

	os.Exit(result.ExitCode)
}

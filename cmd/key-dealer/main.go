package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"key-dealer/pkg/config"
	keys "key-dealer/pkg/key-dealer"

	"github.com/sirupsen/logrus"
)

func main() {

	var configuration config.Configuration

	configPath := "config/configuration.yaml"
	configuration, err := config.LoadConfig(configPath)
	if err != nil {
		logrus.Fatal("Cannot load config")
	}
	level, err := logrus.ParseLevel(configuration.LogLevel)
	if err != nil {
		logrus.Fatal("Log level invalid")
	}

	if configuration.Env != "prod" && configuration.Env != "prd" {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(level)

	// Schedule background task
	go func() {
		for {
			now := time.Now()
			next := now.AddDate(0, 0, 1).Truncate(24 * time.Hour).Add(time.Hour) // Next day at 1:00 AM
			tmr := time.NewTimer(next.Sub(now))
			<-tmr.C
			err := keys.DeleteKeys()
			if err != nil {
				logrus.Errorf("Error when deleting keys: %v", err)
			}
		}
	}()

	logrus.Info(fmt.Sprintf("Serving on http://%s:%s ...", configuration.Hostname, configuration.Port))
	http.HandleFunc("GET /keys/{sa_email}", keys.MakeKeyHandler())

	err = http.ListenAndServe(fmt.Sprintf(":%s", configuration.Port), nil)
	if err != nil {
		logrus.Fatal(err)
	}
}

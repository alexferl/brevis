package main

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/spf13/viper"
	"log/syslog"
)

// InitLogging initializes the logger based on the config
func InitLogging() {
	logFile := viper.GetString("log-file")
	logFormat := viper.GetString("log-format")
	logLevel := viper.GetString("log-level")

	switch logFile {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "null":
		log.SetOutput(ioutil.Discard)
	case "syslog":
		syslogAddr := viper.GetString("syslog-address")
		hook, err := logrus_syslog.NewSyslogHook("udp", syslogAddr, syslog.LOG_INFO, "")
		if err != nil {
			log.Errorf("Unable to connect to syslog server:", err)
		} else {
			log.AddHook(hook)
		}
	default:
		file, err := os.Create(logFile)
		if err != nil {
			log.Warnf("Couldn't open log-file '%s', falling back to stdout: %v", logFile, err)
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(file)
		}

	}

	switch logFormat {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Warnf("Unknown log-format '%s', falling back to 'text' format.", logFormat)
		log.SetFormatter(&log.TextFormatter{})
	}

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.Warnf("Unknown log-level '%s', falling back to 'info' level.", logLevel)
		log.SetLevel(log.InfoLevel)
	}
}

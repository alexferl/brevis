package main

import (
	"github.com/labstack/gommon/log"
	"net"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	Address             net.IP
	Port                uint
	BaseUrl             string
	LogFile             string
	LogFormat           string
	LogLevel            string
	LogRequestsDisabled bool
	BackendType         string
	MongoDBBackend      MongoDBBackend
}

type MongoDBBackend struct {
	Timeout time.Duration
	Uri     string
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	cnf := Config{
		Address:             net.ParseIP("127.0.0.1"),
		Port:                1323,
		BaseUrl:             "http://localhost:1323/",
		LogFile:             "stdout",
		LogFormat:           "text",
		LogLevel:            "info",
		LogRequestsDisabled: false,
		BackendType:         "mongodb",
		MongoDBBackend: MongoDBBackend{
			Timeout: time.Duration(5) * time.Second,
			Uri:     "mongodb://127.0.0.1",
		},
	}
	return &cnf
}

// addFlags adds all the flags from the command line
func (cnf *Config) addFlags(fs *pflag.FlagSet) {
	fs.IPVar(&cnf.Address, "address", cnf.Address, "The IP address to listen at.")
	fs.UintVar(&cnf.Port, "port", cnf.Port, "The port to listen at.")
	fs.StringVar(&cnf.BaseUrl, "base-url", cnf.BaseUrl, "Base URL to prefix short URLs with")
	fs.StringVar(&cnf.LogFile, "log-file", cnf.LogFile, "The log file to write to. "+
		"'stdout' means log to stdout, 'stderr' means log to stderr "+
		"and 'off' means discard log messages.")
	fs.StringVar(&cnf.LogFormat, "log-format", cnf.LogFormat,
		"The log format. Valid format values are: text, json.")
	fs.StringVar(&cnf.LogLevel, "log-level", cnf.LogLevel, "The granularity of log outputs. "+
		"Valid log levels: debug, info, warning, error and critical.")
	fs.BoolVar(&cnf.LogRequestsDisabled, "log-requests-disabled", cnf.LogRequestsDisabled, "Log HTTP requests.")
	fs.StringVar(&cnf.BackendType, "backend-type", cnf.BackendType,
		"Type of backend to use to store short URLs")
	fs.DurationVar(&cnf.MongoDBBackend.Timeout, "backend-mongodb-timeout", cnf.MongoDBBackend.Timeout,
		"Timeout connecting/reading/writing to MongoDB")
	fs.StringVar(&cnf.MongoDBBackend.Uri, "backend-mongodb-uri", cnf.MongoDBBackend.Uri,
		"URI of the MongoDB server")
}

// wordSepNormalizeFunc changes all flags that contain "_" separators
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

// BindFlags normalizes and parses the command line flags
func (cnf *Config) BindFlags() {
	cnf.addFlags(pflag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	pflag.Parse()

	viper.SetEnvPrefix("brevis")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	b := BackendFactory()
	err := b.Init()
	if err != nil {
		log.Fatalf("Error initialising backend: %v", err)
	}

	viper.Set("backend", b)
}

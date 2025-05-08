package config

import (
	"github.com/jessevdk/go-flags"
)

type DatabaseConfiguration struct {
	DSName string `long:"ds" env:"DATASTORE" description:"DataStore name (format: mongo/null)" required:"false" default:"mongo"`
	DSDB   string `long:"ds-db" env:"DATASTORE_DB" description:"DataStore database name" required:"false" default:"iinservice"`
	DSURL  string `long:"ds-url" env:"DATASTORE_URL" description:"DataStore URL" required:"false" default:"mongodb://localhost:27017"`
}

type Config struct {
	Database DatabaseConfiguration
	Port     string `long:"port" env:"PORT" description:"HTTP server port" required:"false" default:"8080"`
}

func Load() (*Config, error) {
	var cfg Config
	parser := flags.NewParser(&cfg, flags.Default)
	if _, err := parser.Parse(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

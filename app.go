package main

import (
	"github.com/clintjedwards/go/config"
	"github.com/clintjedwards/go/storage"
	"go.uber.org/zap"
)

type app struct {
	config  *config.Config
	storage storage.BoltDB
}

func newApp() *app {

	config, err := config.FromEnv()
	if err != nil {
		zap.S().Fatalw("could not load config", "error", err)
	}

	storage, err := storage.NewBoltDB(config.DBPath)
	if err != nil {
		zap.S().Fatalw("could not configure storage", "error", err)
	}

	return &app{
		config:  config,
		storage: storage,
	}
}

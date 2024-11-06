package main

import (
	"SharepointBot/config"
	"SharepointBot/db"
	"go.uber.org/zap"
)

type httpImpl struct {
	logger *zap.SugaredLogger
	db     db.SQL
	config config.Config
}

type HTTP interface {
	// sharepoint.go
	SharepointGoroutine()
}

func NewHTTPInterface(logger *zap.SugaredLogger, db db.SQL, config config.Config) HTTP {
	return &httpImpl{
		logger: logger,
		db:     db,
		config: config,
	}
}

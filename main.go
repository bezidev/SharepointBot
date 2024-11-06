package main

import (
	"SharepointBot/config"
	"SharepointBot/db"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Starting server...")

	var logger *zap.Logger
	var err error

	cfg, err := config.GetConfig()
	if err != nil {
		panic("Error while retrieving config: " + err.Error())
		return
	}

	if cfg.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err.Error())
		return
	}

	sugared := logger.Sugar()

	database, err := db.NewSQL(cfg.DatabaseName, cfg.DatabaseConfig, sugared)
	database.Init()

	if err != nil {
		sugared.Fatal("Error while creating database: ", err.Error())
		return
	}

	sugared.Info("Database created successfully")

	httphandler := NewHTTPInterface(sugared, database, cfg)
	httphandler.SharepointGoroutine()
}

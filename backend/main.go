package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bongofriend/bongo-notes/backend/lib/api"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
	"github.com/bongofriend/bongo-notes/backend/migrations"
	_ "github.com/mattn/go-sqlite3"
)

//	@title		Bongo Notes backend
//	@version	1.0

// @host						localhost:8888
// @BasePath					/
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	appContext, cancel := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := migrations.ApplyMigrations(config); err != nil {
		log.Fatal(err)
	}

	doneCh := make(chan struct{})
	go api.InitApi(appContext, doneCh, config)

	<-signalCh
	cancel()
	<-doneCh
}

func getConfig() (config.Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./local.config.yaml", "Path to config file")
	flag.Parse()
	config, err := config.LoadConfig(configPath)
	if err != nil {
		return config, err
	}
	return config, nil

}

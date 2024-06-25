package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
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
	missingBinary, ok := validateEnvironment()
	if !ok {
		log.Fatalf("%s not found", missingBinary)
	}

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

	errCh := make(chan struct{})
	doneCh := make(chan struct{})
	go api.InitApi(appContext, errCh, doneCh, config)

	defer func() {
		cancel()
		<-doneCh
		close(doneCh)
		close(signalCh)
		close(errCh)
	}()

	select {
	case <-signalCh: //Signal received from OS for termination; cancel and wait for completion
		return
	case <-errCh: //Service encountered error during start up; cancel and wait for completion
		return
	}
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

func validateEnvironment() (string, bool) {
	binaries := []string{"diff", "patch"}
	for _, b := range binaries {
		if !isCommandAvailable(b) {
			return b, false
		}
	}
	return "", true
}

func isCommandAvailable(cmd string) bool {
	res, err := exec.LookPath(cmd)
	log.Println(err)
	return err != nil || res == ""
}

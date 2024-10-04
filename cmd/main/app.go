package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stasdashkevitch/rest-api/cmd/internal/config"
	"github.com/stasdashkevitch/rest-api/cmd/internal/user"
	"github.com/stasdashkevitch/rest-api/cmd/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("--create server")
	router := httprouter.New()

	cfg := config.GetConfig()

	logger.Info("--register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("--start application")

	var listener net.Listener
	var listenError error

	if cfg.Listen.Type == "sock" {
		logger.Info("--detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("--create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("--listen unix socket")
		listener, listenError = net.Listen("unix", socketPath)
		logger.Infof("--server is listeining unix socket %s", socketPath)
	} else {
		logger.Info("--listen tcp")
		listener, listenError = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Infof("--server is listeining port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	}

	if listenError != nil {
		logger.Fatal(listenError)
	}

	server := &http.Server{
		Handler:      router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	logger.Fatal(server.Serve(listener))
}

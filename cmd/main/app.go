package main

import (
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stasdashkevitch/rest-api/cmd/internal/user"
	"github.com/stasdashkevitch/rest-api/cmd/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("--create server")
	router := httprouter.New()

	logger.Info("--register user handler")
	handler := user.NewHandler(*logger)
	handler.Register(router)

	start(router)
}

func start(router *httprouter.Router) {
	logger := logging.GetLogger()
	logger.Info("--start application")
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	logger.Fatal(server.Serve(listener))
}

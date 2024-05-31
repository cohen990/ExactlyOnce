package server

import (
	"context"
	"net/http"

	"github.com/cohen990/exactlyOnce/logging"
)

var logger = logging.NewRoot("server")

type Server struct {
	server *http.Server
}

func Start(port string) Server {
	logger := logger.Child("Start")
	logger.Info("Starting the server on port: %s", port)
	server := &http.Server{Addr: "localhost:" + port, Handler: nil}
	go server.ListenAndServe()
	logger.Info("Server running in background")
	return Server{server}
}

func (server *Server) Shutdown() {
	logger := logger.Child("Shutdown")
	logger.Info("Shutting down the server")
	go server.server.Shutdown(context.Background())
}

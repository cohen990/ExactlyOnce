package server

import (
	"context"
	"net/http"

	"github.com/cohen990/exactlyOnce/logging"
)

var logger = logging.NewRoot("server").Mute()

type Server struct {
	http   *http.Server
	router *http.ServeMux
}

func New(address string) Server {
	server := http.Server{Addr: address}
	mux := http.NewServeMux()
	server.Handler = mux
	return Server{http: &server, router: mux}
}

func (server *Server) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	server.router.HandleFunc(path, handler)
}

func (server *Server) Start() {
	logger := logger.Child("Start")
	logger.Info("Starting the server at %s", server.http.Addr)
	go server.http.ListenAndServe()
	logger.Info("Server running in background")
}

func (server *Server) Shutdown() {
	logger := logger.Child("Shutdown")
	logger.Info("Shutting down the server")
	go server.http.Shutdown(context.Background())
}

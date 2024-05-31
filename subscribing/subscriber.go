package subscribing

import (
	"io"
	"math/rand/v2"

	"net/http"

	"github.com/cohen990/exactlyOnce/logging"
	"github.com/cohen990/exactlyOnce/server"
)

var logger = logging.NewRoot("subscriber")

type Subscriber struct {
	logRoot            logging.LogRoot
	ReceivedCount      int
	ReceiveFailedCount int
	Url                string
	server             server.Server
	port               string
}

func (subscriber *Subscriber) process(message string) Status {
	log := logger.Child("Receive")
	if rand.Float32() > 0.5 {
		log.Info("received message: %q", message)
		subscriber.ReceivedCount++
		return Received
	} else {
		log.Info("Failed to process message: %q", message)
		subscriber.ReceiveFailedCount++
		return Failed
	}
}

func (subscriber *Subscriber) ReceiveHttp(response http.ResponseWriter, request *http.Request) {
	buffer, err := io.ReadAll(request.Body)
	if err != nil {
		panic("oh fuck")
	}

	result := subscriber.process(string(buffer))
	if result == Received {
		response.WriteHeader(http.StatusOK)
	} else {
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (subscriber *Subscriber) Initialise() {
	http.HandleFunc("/receive", subscriber.ReceiveHttp)
	// httptest for testing
	// checkfunc
	subscriber.port = "8082"
	subscriber.Url = "http://localhost:" + subscriber.port
	subscriber.logRoot = logging.NewRoot("subscriber")
}

func (subscriber *Subscriber) Start() {
	subscriber.server = server.Start(subscriber.port)
}

func (subscriber *Subscriber) Shutdown() {
	subscriber.server.Shutdown()
}

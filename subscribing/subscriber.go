package subscribing

import (
	"io"
	"math/rand/v2"

	"net/http"

	"github.com/cohen990/exactlyOnce/logging"
)

var logger = logging.Local("subscriber")

type Subscriber struct {
	ReceivedCount      int
	ReceiveFailedCount int
	PanickedCount      int
	Url                string
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

func (subscriber *Subscriber) Start() {
	logger := logger.Child("Start")
	port := "8081"
	logger.Info("Starting the server on port: %s", port)
	subscriber.Url = "http://localhost:" + port
	http.HandleFunc("/receive", subscriber.ReceiveHttp)
	go http.ListenAndServe("localhost:"+port, nil)
	logger.Info("Server running in background")
}

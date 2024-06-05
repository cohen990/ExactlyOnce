package brokering

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"

	"github.com/cohen990/exactlyOnce/logging"
	"github.com/cohen990/exactlyOnce/server"

	"github.com/google/uuid"
)

type Broker struct {
	logRoot           logging.LogRoot
	subscriberUrl     string
	messageQueue      map[uuid.UUID]string
	requeuedMessages  map[uuid.UUID]string
	Brokering         bool
	TotalSentCount    int
	SentCount         int
	SendFailedCount   int
	SendPanickedCount int
	Url               string
	server            server.Server
}

func (broker *Broker) Initialise() {
	broker.messageQueue = make(map[uuid.UUID]string)
	broker.requeuedMessages = make(map[uuid.UUID]string)
	broker.Url = "http://localhost:8082"
	broker.server = server.New(broker.Url)
	broker.server.HandleFunc("/enqueue", broker.EnqueueHttp)
	broker.logRoot = logging.NewRoot("subscriber").Mute()
}
func (broker *Broker) RegisterSubscriber(url string) {
	broker.subscriberUrl = url
}

func (broker *Broker) enqueue(message string) EnqueuedStatus {
	if rand.Float32() > 0.5 {
		broker.messageQueue[uuid.New()] = message
		return Enqueued
	} else {
		return EnqueuingFailed
	}
}

func (broker *Broker) Send(message string, status chan SendStatus) {
	logger := broker.logRoot.Child("Send")

	logger.Info("Sending %q", message)
	response, err := http.Post(broker.subscriberUrl+"/receive", "text/plain", bytes.NewBuffer([]byte(message)))
	if err != nil {
		logger.Info("Error when sending message %q: %q", message, err)
		status <- SendFailed
		return
	}

	broker.TotalSentCount++
	if response.StatusCode == 200 {
		status <- Sent
	} else {
		status <- SendFailed
	}
	logger.Info("Message %q sent. Status Code: %s", message, response.StatusCode)
}

func (broker *Broker) Process() {
	logger := broker.logRoot.Child("Process")
	logger.Info("Brokering %d messages.", len(broker.messageQueue))

	broker.Brokering = true
	clear(broker.requeuedMessages)

	for {
		logger.Info("Queue size: %d", len(broker.messageQueue))
		sendStatus := make(chan SendStatus)
		if len(broker.messageQueue) == 0 {
			break
		}
		for id, message := range broker.messageQueue {
			logger.Info("Initiating send of %q", message)
			go broker.Send(message, sendStatus)
			result := <-sendStatus
			logger.Info("Send of %q completed", message)
			if result == Sent {
				broker.SentCount++
			} else {
				broker.SendFailedCount++
				logger.Info("requeueing %q", message)
				broker.requeuedMessages[id] = message
			}
		}
		logger.Info("Requeueing %d failed messages", len(broker.requeuedMessages))
		broker.requeue()
	}
}

func (broker *Broker) requeue() {
	broker.requeuedMessages, broker.messageQueue = broker.messageQueue, broker.requeuedMessages
	broker.requeuedMessages = make(map[uuid.UUID]string)
}

func (broker *Broker) Start() {
	broker.server.Start()
}

func (broker *Broker) Shutdown() {
	broker.server.Shutdown()
}

func (broker *Broker) EnqueueHttp(response http.ResponseWriter, request *http.Request) {
	logger := broker.logRoot.Child("EnqueueHttp")

	buffer, err := io.ReadAll(request.Body)
	if err != nil {
		panic("oh fuck")
	}

	result := broker.enqueue(string(buffer))
	if result == Enqueued {
		logger.Info("Enqueued successfully")
		response.WriteHeader(http.StatusOK)
	} else {
		logger.Info("Failed to enqueue")
		response.WriteHeader(http.StatusInternalServerError)
	}
}

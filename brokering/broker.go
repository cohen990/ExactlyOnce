package brokering

import (
	"reflect"

	"github.com/cohen990/exactlyOnce/logging"

	"github.com/cohen990/exactlyOnce/subscribing"

	"github.com/google/uuid"
)

var logger = logging.Local("broker")

type Broker struct {
	subscriber        *subscribing.Subscriber
	messageQueue      map[uuid.UUID]string
	requeuedMessages  map[uuid.UUID]string
	Brokering         bool
	SentCount         int
	SendFailedCount   int
	SendPanickedCount int
}

func (broker *Broker) Initialise() {
	broker.messageQueue = make(map[uuid.UUID]string)
	broker.requeuedMessages = make(map[uuid.UUID]string)
}
func (broker *Broker) Register(subscriber *subscribing.Subscriber) {
	broker.subscriber = subscriber
}

func (broker *Broker) Enqueue(message string) EnqueuedStatus {
	broker.messageQueue[uuid.New()] = message
	return Enqueued
}

func (broker *Broker) Send(message string, sendStatus chan SendStatus) {
	logger := logger.Child("Send")
	defer func() {
		if r := recover(); r != nil {
			i := 1
			_ = i
			// no op - stop compiler whining
		}
	}()

	logger.Info("Sending %q", message)
	status := make(chan subscribing.Status)
	go broker.subscriber.Receive(message, status)
	result := <-status
	logger.Info("Message %q sent. Status: %s", message, result)
	if result == subscribing.Failed {
		sendStatus <- SendFailed
	}
	sendStatus <- Sent
}

func (broker *Broker) Process() {
	logger := logger.Child("Process")
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
	requeuePointer := reflect.ValueOf(&broker.requeuedMessages).Elem()
	queuePointer := reflect.ValueOf(&broker.messageQueue).Elem()
	queuePointer.Set(requeuePointer)
	broker.requeuedMessages = make(map[uuid.UUID]string)
}

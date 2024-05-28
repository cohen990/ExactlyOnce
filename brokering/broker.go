package brokering

import (
	"fmt"

	"github.com/cohen990/exactlyOnce/subscribing"
)

type Broker struct {
	subscriber        *subscribing.Subscriber
	messageQueue      []string
	requeuedMessages  []string
	Brokering         bool
	SentCount         int
	SendFailedCount   int
	SendPanickedCount int
}

func (broker *Broker) Register(subscriber *subscribing.Subscriber) {
	broker.subscriber = subscriber
}

func (broker *Broker) Enqueue(message string) EnqueuedStatus {
	broker.messageQueue = append(broker.messageQueue, message)
	return Enqueued
}

func (broker *Broker) Send(message string) SendStatus {
	defer func() {
		if r := recover(); r != nil {
			i := 1
			_ = i
			// no op - stop compiler whining
		}
	}()

	result := broker.subscriber.Receive(message)
	if result == subscribing.Failed {
		return SendFailed
	}
	return Sent
}

func (broker *Broker) Process() {
	fmt.Printf("Brokering %d messages.\n", len(broker.messageQueue))

	broker.Brokering = true
	broker.requeuedMessages = []string{}

	for {
		if len(broker.messageQueue) == 0 {
			break
		}
		for _, message := range broker.messageQueue {
			result := broker.Send(message)
			if result == Sent {
				broker.SentCount++
			} else {
				broker.SendFailedCount++
				broker.requeuedMessages = append(broker.requeuedMessages, message)
			}
		}
		broker.messageQueue = broker.requeuedMessages
		broker.requeuedMessages = []string{}
	}
}

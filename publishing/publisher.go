package publishing

import (
	"github.com/cohen990/exactlyOnce/brokering"
)

type Publisher struct {
	EnqueuedCount      int
	EnqueueFailedCount int
}

func (publisher *Publisher) Publish(broker *brokering.Broker, message string) {
	result := broker.Enqueue(message)
	if result == brokering.Enqueued {
		publisher.EnqueuedCount++
	} else {
		publisher.EnqueueFailedCount++
	}
}

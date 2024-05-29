package subscribing

import (
	"math/rand/v2"

	"github.com/cohen990/exactlyOnce/logging"
)

var logger = logging.Local("subscriber")

type Subscriber struct {
	ReceivedCount      int
	ReceiveFailedCount int
	PanickedCount      int
}

func (subscriber *Subscriber) Receive(message string, status chan Status) {
	log := logger.Child("Receive")
	if rand.Float32() > 0.5 {
		log.Info("received message: %q", message)
		subscriber.ReceivedCount++
		status <- Received
		// } else if rand.Float32() > 0.5 {
		// 	subscriber.PanickedCount++
		// 	panic("explode")
	} else {
		log.Info("Failed to process message: %q", message)
		subscriber.ReceiveFailedCount++
		status <- Failed
	}
}

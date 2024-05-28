package subscribing

import (
	"fmt"
	"math/rand/v2"
)

type Subscriber struct {
	ServiceName        string
	ReceivedCount      int
	ReceiveFailedCount int
	PanickedCount      int
}

func (subscriber *Subscriber) Receive(message string) Status {
	if rand.Float32() > 0.5 {
		fmt.Printf("%s received message: %q\n", subscriber.ServiceName, message)
		subscriber.ReceivedCount++
		return Received
	} else if rand.Float32() > 0.5 {
		subscriber.PanickedCount++
		panic("explode")
	} else {
		subscriber.ReceiveFailedCount++
		return Failed
	}
}

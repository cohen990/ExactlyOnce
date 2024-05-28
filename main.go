package main

import (
	"fmt"

	"github.com/cohen990/exactlyOnce/brokering"
	"github.com/cohen990/exactlyOnce/publishing"
	"github.com/cohen990/exactlyOnce/subscribing"
)

func main() {
	broker := brokering.Broker{}
	subscriber := subscribing.Subscriber{ServiceName: "messageSubscriberService"}
	publisher := publishing.Publisher{}
	broker.Register(&subscriber)
	messages := []string{"Hi", "bye", "why"}

	for _, message := range messages {
		publisher.Publish(&broker, message)
	}

	broker.Process()

	fmt.Printf("Publishing %d messages, completed.\n", len(messages))
	fmt.Printf("Publisher failed to enqueue messages %d times.\n", publisher.EnqueueFailedCount)
	fmt.Printf("Publisher enqueued %d messages.\n", publisher.EnqueuedCount)
	fmt.Printf("Broker failed to send messages %d times.\n", broker.SendFailedCount)
	fmt.Printf("Broker sent %d messages.\n", broker.SentCount)
	fmt.Printf("Subscriber panicked %d times.\n", subscriber.PanickedCount)
	fmt.Printf("Subscriber failed to receive messages %d times.\n", subscriber.ReceiveFailedCount)
	fmt.Printf("Subscriber received %d messages.\n", subscriber.ReceivedCount)

	fmt.Printf("\nPublished %d messages. Received %d messages.\n", len(messages), subscriber.ReceivedCount)
	fmt.Printf("\nExactly once delivery achieved?: %t\n", len(messages) == subscriber.ReceivedCount)
}

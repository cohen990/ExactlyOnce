package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/cohen990/exactlyOnce/brokering"
	"github.com/cohen990/exactlyOnce/logging"
	"github.com/cohen990/exactlyOnce/publishing"
	"github.com/cohen990/exactlyOnce/subscribing"
	"github.com/tjarratt/babble"
)

func main() {
	logger := logging.Local("main")
	broker := brokering.Broker{}
	subscriber := subscribing.Subscriber{}
	publisher := publishing.Publisher{}
	subscriber.Start()
	broker.Initialise()
	broker.Register(&subscriber)
	brokerServer := broker.Start()

	// messageCount := rand.IntN(30)
	messageCount := 3
	babbler := babble.NewBabbler()
	babbler.Count = messageCount
	babbler.Separator = ","
	messages := strings.Split(babbler.Babble(), ",")
	retries := []string{}

	for len(messages) > 0 {
		logger.Info("Publishing %d messages", len(messages))
		for _, message := range messages {
			var result publishing.PublishStatus
			if rand.Float32() > 0.5 {
				logger.Info("Shutting down the server")
				brokerServer.Shutdown(context.Background())
				result = publisher.Publish(broker.Url, message)
				logger.Info("Restoring the server")
				brokerServer = broker.Start()
			} else {
				result = publisher.Publish(broker.Url, message)
			}
			if result != publishing.Published {
				retries = append(retries, message)
			}
		}
		time.Sleep(100000)
		logger.Info("Retrying %d messages", len(retries))
		messages = retries
		retries = []string{}
		logger.Info("Retries reset: %d", len(retries))
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

	fmt.Printf("\nQueued %d messages. Sent %d messages in total.\n", len(messages), broker.TotalSentCount)
	fmt.Printf("\nExactly once delivery achieved?: %t\n", len(messages) == broker.TotalSentCount)
}

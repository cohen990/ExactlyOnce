package main

import (
	"context"
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

	logger.Info("============================")
	logger.Info("Publishing %d messages, completed.", messageCount)
	logger.Info("Publisher failed to enqueue messages %d times.", publisher.EnqueueFailedCount)
	logger.Info("Publisher enqueued %d messages.", publisher.EnqueuedCount)
	logger.Info("Broker failed to send messages %d times.", broker.SendFailedCount)
	logger.Info("Broker sent %d messages.", broker.SentCount)
	logger.Info("Subscriber panicked %d times.", subscriber.PanickedCount)
	logger.Info("Subscriber failed to receive messages %d times.", subscriber.ReceiveFailedCount)
	logger.Info("Subscriber received %d messages.", subscriber.ReceivedCount)

	logger.Info("Queued %d messages. Sent %d messages in total.", messageCount, broker.TotalSentCount)
	logger.Info("Exactly once delivery achieved?: %t", messageCount == broker.TotalSentCount)
}

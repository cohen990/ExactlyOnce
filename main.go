package main

import (
	"math/rand/v2"
	"strings"

	"github.com/cohen990/exactlyOnce/brokering"
	"github.com/cohen990/exactlyOnce/chaos"
	"github.com/cohen990/exactlyOnce/logging"
	"github.com/cohen990/exactlyOnce/publishing"
	"github.com/cohen990/exactlyOnce/subscribing"
	"github.com/tjarratt/babble"
)

func main() {
	logger := logging.NewRoot("exactlyOnce").Child("main")
	broker := brokering.Broker{}
	subscriber := subscribing.Subscriber{}
	publisher := publishing.Publisher{}
	subscriber.Initialise()
	subscriber.Start()
	broker.Initialise()
	publisher.Initialise()
	broker.Start()
	broker.RegisterSubscriber(subscriber.Url)

	outage := brokering.NewBrokerOutage(&broker)
	chaos := chaos.Chaos[chaos.Action[publishing.PublishStatus], publishing.PublishStatus]{}
	chaos.Register(outage)

	messageCount := rand.IntN(30)
	// messageCount := 3
	babbler := babble.NewBabbler()
	babbler.Count = messageCount
	babbler.Separator = ","
	messages := strings.Split(babbler.Babble(), ",")
	retries := []string{}

	for len(messages) > 0 {
		logger.Info("=============New run=============")
		logger.Info("Publishing %d messages", len(messages))
		for _, message := range messages {
			result := chaos.InjectChaos(func() publishing.PublishStatus { return publisher.Publish(broker.Url, message) })
			if result != publishing.Published {
				retries = append(retries, message)
			}
		}
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
	logger.Info("Subscriber failed to receive messages %d times.", subscriber.ReceiveFailedCount)
	logger.Info("Subscriber received %d messages.", subscriber.ReceivedCount)

	logger.Info("Queued %d messages. Sent %d messages in total.", messageCount, broker.TotalSentCount)
	logger.Info("Exactly once delivery achieved?: %t", messageCount == broker.TotalSentCount)
}

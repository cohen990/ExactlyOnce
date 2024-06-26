package publishing

import (
	"bytes"
	"net/http"

	"github.com/cohen990/exactlyOnce/logging"
)

type Publisher struct {
	logRoot            logging.LogRoot
	EnqueuedCount      int
	EnqueueFailedCount int
}

func (publisher *Publisher) Initialise() {
	publisher.logRoot = logging.NewRoot("publisher")
}

func (publisher *Publisher) Publish(brokerUrl string, message string) PublishStatus {
	logger := publisher.logRoot.Child("Publish")
	logger.Info("Publishing %q", message)
	response, err := http.Post(brokerUrl+"/enqueue", "text/plain", bytes.NewBuffer([]byte(message)))
	if err != nil {
		logger.Info("Error when posting %q", message)
		publisher.EnqueueFailedCount++
		return Failed
	}

	if response.StatusCode == 200 {
		logger.Info("Successfully published %q", message)
		publisher.EnqueuedCount++
		return Published
	} else {
		logger.Info("Failed to publish %q", message)
		publisher.EnqueueFailedCount++
		return Failed
	}
}

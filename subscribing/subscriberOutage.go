package subscribing

import (
	"github.com/cohen990/exactlyOnce/logging"
)

type SubscriberOutage struct {
	logRoot    logging.LogRoot
	subscriber *Subscriber
}

func MakeSubscriberOutage(subscriber *Subscriber) SubscriberOutage {
	return SubscriberOutage{
		subscriber: subscriber,
		logRoot:    logging.NewRoot("susbcriberOutage")}
}

func (outage SubscriberOutage) Break() {
	logger := logger.Child("Break")
	outage.subscriber.Shutdown()
	logger.Info("The broker has been taken down")
}

func (outage SubscriberOutage) Fix() {
	logger := logger.Child("Fix")
	outage.subscriber.Start()
	logger.Info("The broker has been restored")
}

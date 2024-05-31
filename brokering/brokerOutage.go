package brokering

import (
	"github.com/cohen990/exactlyOnce/logging"
)

var logger = logging.NewRoot("brokerOutage")

type BrokerOutage struct {
	Broker  *Broker
	logRoot logging.LogRoot
}

func NewBrokerOutage(broker *Broker) BrokerOutage {
	return BrokerOutage{Broker: broker, logRoot: logging.NewRoot("brokerOutage")}
}

func (outage BrokerOutage) Break() {
	logger := logger.Child("Break")
	outage.Broker.Shutdown()
	logger.Info("The broker has been taken down")
}

func (outage BrokerOutage) Fix() {
	logger := logger.Child("Fix")
	outage.Broker.Start()
	logger.Info("The broker has been restored")
}

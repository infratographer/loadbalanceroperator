// Package events provides functions required for responding to
// various events and connecting with the event queue
package events

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// ConnectNATS establishes a connection with a specified NATS server
func ConnectNATS(uri string, logger *zap.SugaredLogger) *nats.Conn {
	logger.Debugf("Connecting to NATS at %s", uri)
	nc, err := nats.Connect(uri)

	if err != nil {
		logger.Fatalf("Unable to connect to NATS at %s", uri)
	}

	return nc
}

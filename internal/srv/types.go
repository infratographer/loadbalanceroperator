package srv

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
)

// Server holds options for server connectivity and settings
type Server struct {
	Context         context.Context
	Logger          *zap.SugaredLogger
	KubeClient      *rest.Config
	JetstreamClient nats.JetStreamContext
	Debug           bool
	Metro           string
	Prefix          string
	ChartPath       string
}

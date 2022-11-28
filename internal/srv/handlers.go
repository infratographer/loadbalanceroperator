package srv

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Moby/Moby/pkg/namesgenerator"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// MessageHandler handles the routing of events from specified queues
func (s *Server) MessageHandler(m *nats.Msg) {
	switch m.Subject {
	case fmt.Sprintf("%s.create", s.Prefix):
		err := s.createMessageHandler(m)
		if err != nil {
			// TODO: eventually we'll want to requeue failed events
			s.Logger.Errorln("unable to process create")
		}
	case fmt.Sprintf("%s.update", s.Prefix):
		err := s.updateMessageHandler(m)
		if err != nil {
			// TODO: eventually we'll want to requeue failed events
			s.Logger.Errorln("unable to process update")
		}
	default:
		s.Logger.Debug("This is some other set of queues that we don't know about.")
	}
}

func (s *Server) createMessageHandler(m *nats.Msg) error {
	name := strings.ReplaceAll(namesgenerator.GetRandomName(0), "_", "-")
	err := s.CreateNamespace(string(m.Data))

	if err != nil {
		return err
	}

	err = s.CreateApp(name, s.ChartPath, string(m.Data))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) updateMessageHandler(m *nats.Msg) error {
	s.Logger.Infoln("updating")
	return nil
}

// ExposeEndpoint exposes a specified port for various checks
func ExposeEndpoint(name string, port string, logger *zap.SugaredLogger) {
	if port == "" {
		logger.Fatalf("port has not been provided for endpoint: %s", name)
	}

	logger.Infof("Starting %s endpoint", name)

	go func() {
		_ = http.ListenAndServe(port, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("ok"))
			},
		))
	}()
}

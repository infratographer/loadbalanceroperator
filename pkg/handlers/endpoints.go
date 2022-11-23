package handlers

import (
	"net/http"

	"go.uber.org/zap"
)

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

package deployments

import (
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func KubeAuth(logger *zap.SugaredLogger, path string) *rest.Config {

	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Debugln("Unable to read in-cluster config")
		if path != "" {
			config, err = clientcmd.BuildConfigFromFlags("", path)
			if err != nil {
				return nil
			}
		} else {
			return nil
		}
	}

	return config
}

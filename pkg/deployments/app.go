package deployments

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func CreateNamespace(client *rest.Config, projectId string, logger *zap.SugaredLogger) error {
	kc, err := kubernetes.NewForConfig(client)
	if err != nil {
		logger.Fatalln("Unable to authenticate against kubernetes cluster")
		return nil
	}

	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: projectId}}
	_, err = kc.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
	return err
}

func CreateApp(name string, client *rest.Config, chartPath string, namespace string, logger *zap.SugaredLogger) error {
	releaseName := fmt.Sprintf("%s-%s", name, namespace)
	releaseNS := namespace

	chart, err := loader.Load(chartPath)
	if err != nil {
		logger.Errorf("Unable to load chart from %s", chartPath)
		return err
	}

	config := new(action.Configuration)
	cliopt := genericclioptions.NewConfigFlags(false)
	wrapper := func(*rest.Config) *rest.Config {
		return client
	}
	cliopt.WithWrapConfigFn(wrapper)

	config.Init(cliopt, namespace, "secret", func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	})

	hc := action.NewInstall(config)
	hc.ReleaseName = releaseName
	hc.Namespace = releaseNS
	_, err = hc.Run(chart, nil)
	if err != nil {
		logger.Errorf("Unable to deploy %s to %s", releaseName, releaseNS)
	}

	return nil
}

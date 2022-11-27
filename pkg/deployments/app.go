// Package deployments handles the necessary functions to deploy an application
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

// CreateNamespace creates a namespace based upon the projectId that is passed
func CreateNamespace(client *rest.Config, projectID string, logger *zap.SugaredLogger) error {
	kc, err := kubernetes.NewForConfig(client)
	if err != nil {
		logger.Fatalln("Unable to authenticate against kubernetes cluster")
		return nil
	}

	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: projectID}}
	_, err = kc.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})

	return err
}

// CreateApp loads our specified helm chart and then deploys the chart
// to the specified namespace
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

	err = config.Init(cliopt, namespace, "secret", func(format string, v ...interface{}) {
		fmt.Println(v)
	})
	if err != nil {
		logger.Errorln("Unable to initialize helm client: %s", err)
	}

	hc := action.NewInstall(config)
	hc.ReleaseName = releaseName
	hc.Namespace = releaseNS
	_, err = hc.Run(chart, nil)

	if err != nil {
		logger.Errorf("Unable to deploy %s to %s", releaseName, releaseNS)
	}

	return nil
}

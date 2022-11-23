/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Moby/Moby/pkg/namesgenerator"
	"github.com/infratographer/wallenda/pkg/deployments"
	"github.com/infratographer/wallenda/pkg/events"
	"github.com/infratographer/wallenda/pkg/handlers"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// processCmd represents the base command when called without any subcommands
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Message machine go brrrrrrrr",
	Long:  `Begin processing requests from queues.`,
	Run: func(cmd *cobra.Command, args []string) {

		liveness_port := viper.GetString("liveness-port")
		handlers.ExposeEndpoint("healthz", liveness_port, logger)

		nats_url := viper.GetString("nats.url")
		nc := events.ConnectNATS(nats_url, logger)
		defer nc.Close()

		js, err := nc.JetStream()
		if err != nil {
			logger.Fatalf("Unable to establish a Jetstream context: %s", err)
		}

		readiness_port := viper.GetString("readiness-port")
		handlers.ExposeEndpoint("readyz", readiness_port, logger)

		subject_prefix := viper.GetString("nats.subject-prefix")
		if subject_prefix == "" {
			logger.Fatalln("NATS subject prefix is not set.")
		}

		stream_name := viper.GetString("nats.stream-name")
		if stream_name == "" {
			logger.Fatalln("NATS stream name is not set.")
		}

		chart := viper.GetString("chart-path")
		if chart == "" {
			logger.Fatalln("No chart was provided.")
		}

		kubeconfig := viper.GetString("kube-config-path")
		client := deployments.KubeAuth(logger, kubeconfig)

		js.QueueSubscribe(fmt.Sprintf("%s.>", subject_prefix), "wallenda-workers", func(m *nats.Msg) {
			fmt.Printf("Msg received on [%s] : %s\n", m.Subject, string(m.Data))
			switch m.Subject {
			case fmt.Sprintf("%s.create", subject_prefix):
				deployments.CreateNamespace(client, string(m.Data), logger)
				//TODO: Just using random name generator for now. This should go away ASAP.
				name := namesgenerator.GetRandomName(0)
				name = strings.Replace(name, "_", "-", -1)
				deployments.CreateApp(name, client, chart, string(m.Data), logger)
			case fmt.Sprintf("%s.update", subject_prefix):
				fmt.Println("zap")
			default:
				logger.Debug("This is some other set of queues that we don't know about.")
			}
		}, nats.BindStream(stream_name))

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGABRT)

		// Wait for appropriate signal to trigger clean shutdown
		for {
			select {
			case recvSig := <-sigCh:
				signal.Stop(sigCh)
				fmt.Printf("\nGet the brooms: %s\n", recvSig)
				return
			}
		}
	},
}

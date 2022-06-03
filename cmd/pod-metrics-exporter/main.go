package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"os"
)

// =========================================================================
// App configuration

type config struct {
	labelName         string
	labelValue        string
	metricsListenAddr string
	kubeconfig        string
	podLabelSelector  string
}

type application struct {
	config          config
	logger          *logrus.Logger
	kubeClient      *kubernetes.Clientset
	podWatchChannel <-chan watch.Event
	podMetric       *prometheus.GaugeVec
}

func main() {

	// =========================================================================
	// Init app configuration

	var cfg config

	var log = logrus.New()
	log.Out = os.Stdout

	app := &application{
		config: cfg,
		logger: log,
	}

	// =========================================================================
	// Setup CLI flags and config

	cliApp := &cli.App{
		Name:  "pod-metrics-exporter",
		Usage: "fight the loneliness!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "label-name",
				Value:       "app",
				Usage:       "name of a pod label to filter on",
				Destination: &app.config.labelName,
			},
			&cli.StringFlag{
				Name:        "label-value",
				Value:       "demo-app",
				Usage:       "value corresponding to the pod label named '--label-name'",
				Destination: &app.config.labelValue,
			},
			&cli.StringFlag{
				Name:        "metrics-listen-addr",
				Value:       ":8080",
				Usage:       "listen address of the prometheus metrics server",
				Destination: &app.config.metricsListenAddr,
			},
			&cli.StringFlag{
				Name:        "kubeconfig",
				Required:    true,
				Usage:       "absolute path to a kubeconfig",
				Destination: &app.config.kubeconfig,
			},
		},
		Action: func(c *cli.Context) error {
			err := app.init(c)
			if err != nil {
				return fmt.Errorf("initializing app: %w", err)
			}

			return app.serve(c)
		},
	}

	// =========================================================================
	// Run the app

	err := cliApp.Run(os.Args)
	if err != nil {
		app.logger.Fatal(err)
	}
}

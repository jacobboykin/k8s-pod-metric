package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func (app *application) init(c *cli.Context) error {
	ctx := c.Context
	app.logger.Info("Starting up application")

	// =========================================================================
	// Configure connection to Kubernetes API server

	var err error
	app.kubeClient, err = app.getKubeClient()
	if err != nil {
		return err
	}

	// =========================================================================
	// Get Pod Resource Version for watching resources

	app.config.podLabelSelector = fmt.Sprintf("%s=%s", app.config.labelName, app.config.labelValue)
	pods, err := app.getPods(ctx, app.kubeClient, metav1.ListOptions{LabelSelector: app.config.podLabelSelector})
	if err != nil {
		return err
	}
	resourceVersion := pods.ListMeta.ResourceVersion

	// =========================================================================
	// Set up a watch channel to catch changes to any Pod Resources

	coreV1Pods := app.kubeClient.CoreV1().Pods("")
	watcher, err := coreV1Pods.Watch(ctx, metav1.ListOptions{ResourceVersion: resourceVersion})
	if err != nil {
		return err
	}
	app.podWatchChannel = watcher.ResultChan()

	// =========================================================================
	// Create the pod count Prometheus gauge metric and update it

	app.podMetric = app.createPodCountMetric()
	err = app.updatePodCountMetric(pods, app.podMetric)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) serve(c *cli.Context) error {

	// =========================================================================
	// Listen to the pod watch channel and update the pod count metric whenever an
	// event is detected.

	ctx := c.Context
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		for {
			app.logger.Infof("Watching for changes to Pod resources")
			event := <-app.podWatchChannel

			if event.Type != "" {
				app.logger.Infof("Pod Event detected: %s", event.Type)
				pods, err := app.getPods(ctx, app.kubeClient, metav1.ListOptions{
					LabelSelector: app.config.podLabelSelector})
				if err != nil {
					return fmt.Errorf("fetching pods: %w", err)
				}

				err = app.updatePodCountMetric(pods, app.podMetric)
				if err != nil {
					return fmt.Errorf("udpating pod_count metric: %w", err)
				}
			}
		}
	})

	// =========================================================================
	// Configure prometheus metrics endpoint

	http.Handle("/metrics", promhttp.Handler())

	// =========================================================================
	// Start server

	eg.Go(func() error {
		app.logger.Infof("Starting HTTP server at http://localhost%s", app.config.metricsListenAddr)
		return http.ListenAndServe(app.config.metricsListenAddr, nil)
	})

	return eg.Wait()
}

package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	v1 "k8s.io/api/core/v1"
)

var (
	labelName  = "label_name"
	labelValue = "label_value"
	phase      = "phase"
)

// createPodCountMetric creates the 'pod_count' prometheus gauge metric.
func (app *application) createPodCountMetric() *prometheus.GaugeVec {
	return promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pod_count",
			Help: "Current number of pods",
		},
		[]string{labelName, labelValue, phase},
	)
}

// updatePodCountMetric updates the 'pod_count' prometheus gauge metric with the information
// from the given PodList.
func (app *application) updatePodCountMetric(pods *v1.PodList, podCountGauge *prometheus.GaugeVec) error {
	var phaseValue string

	podsCount := len(pods.Items)

	// Multiple Pods could have different statuses. To implement the functional requirement,
	// just fetch the Status Phase of the first Pod in the list.

	if podsCount > 0 {
		phaseValue = string(pods.Items[0].Status.Phase)
	}

	podCountGauge.With(prometheus.Labels{
		labelName:  app.config.labelName,
		labelValue: app.config.labelValue,
		phase:      phaseValue,
	}).Set(float64(podsCount))

	app.logger.Infof("Updated pod_count metric")

	return nil
}

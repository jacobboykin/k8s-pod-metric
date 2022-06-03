package main

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// getKubeClient creates a kubernetes.Clientset from the absolute path of a kubeconfig file.
func (app *application) getKubeClient() (*kubernetes.Clientset, error) {
	app.logger.Infof("Configuring connection to Kubernetes Cluster using kubeconfig file: %s",
		app.config.kubeconfig)
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", app.config.kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

// getPods fetches Pods from all namespaces in the cluster and includes the usage of ListOptions.
func (app *application) getPods(ctx context.Context, kubeClient *kubernetes.Clientset,
	listOptions metav1.ListOptions) (*v1.PodList, error) {

	coreV1Pods := kubeClient.CoreV1().Pods("")
	pods, err := coreV1Pods.List(ctx, listOptions)
	if err != nil {
		return nil, err
	}

	return pods, nil
}

SHELL := /bin/bash

# =========================================================================
# Config

WORKING_DIR := ${CURDIR}
KUBECONFIG_PATH := "${WORKING_DIR}/kubeconfig"

# =========================================================================
# Run

# Build and run the application via `go run` using the default flag values
run-go:
	go run ./cmd/pod-metrics-exporter --kubeconfig ${KUBECONFIG_PATH}

# =========================================================================
# Bootstrap

# Create k8s cluster and create kubeconfig file
kube-create-cluster:
	k3d cluster create planet-dev-test
	k3d kubeconfig get planet-dev-test > kubeconfig

kube-cleanup:
	k3d cluster delete planet-dev-test
	rm -f ${KUBECONFIG_PATH}

# Deploy test resources to k8s cluster
kube-deploy-test-resources:
	kubectl --kubeconfig ${KUBECONFIG_PATH} apply -f ./k8s

# Fetch test pods with label selector used in default pod-metrics-exporter args
kube-get-demo-app-pods:
	kubectl --kubeconfig ${KUBECONFIG_PATH} get po -l app=demo-app -A

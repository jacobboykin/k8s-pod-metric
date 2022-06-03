# Pod Metrics Exporter

The Pod Metrics Exporter is a Go command line application that:
* maintains a custom `pod_count` [Prometheus Gauge Metric](https://prometheus.io/docs/concepts/metric_types/) with the total number of Kubernetes Pods that match a given [Label Selector](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors).
* runs an HTTP server to host the metrics at `http://<hostname>:<port>/metrics`

```shell
$ go run ./cmd/pod-metrics-exporter \
  --label-name app \
  --label-value demo-app

$ curl -s localhost:8080/metrics | grep pod_count
# HELP pod_count Current number of pods
# TYPE pod_count gauge
pod_count{label_name="app",label_value="demo-app",phase="Running"} 2
```

## Running the application locally
### Requirements
* macOS/Linux
* Go (v1.17+)
* Docker
* [k3d](https://k3d.io/v5.4.1/)
* kubectl

1. Stand up a test k8s cluster
    ```shell
    $ make kube-create-cluster 
    ```
2. Deploy the test k8s resources in [k8s/](k8s)
   ```shell
   $ make kube-deploy-test-resources
   ```
3. Run the application with `go run`
   ```shell
    $ make run-go 
      go run ./cmd/pod-metrics-exporter --kubeconfig "/Users/jacob/dev/github.com/jacobboykin/hobbes-pod-metrics-test/kubeconfig"
      INFO[0000] Starting up application                      
      INFO[0000] Configuring connection to Kubernetes Cluster using kubeconfig file: /Users/jacob/dev/github.com/jacobboykin/hobbes-pod-metrics-test/kubeconfig 
      INFO[0000] Fetching pods from all namespaces using label selector: app=demo-app 
      INFO[0000] Updated pod_count metric                     
      INFO[0000] Starting HTTP server at http://localhost:8080 
      INFO[0000] Watching for changes to Pod resources 
   ```
4. Query the 'pod_count' metric and validate the response
   ```shell
   $ make kube-get-demo-app-pods
   kubectl --kubeconfig "/Users/jacob/dev/github.com/jacobboykin/hobbes-pod-metrics-test/kubeconfig" get po -l app=demo-app -A
   NAMESPACE         NAME       READY   STATUS    RESTARTS   AGE
   demo-staging      demo-app   1/1     Running   0          8m2s
   demo-production   demo-app   1/1     Running   0          8m2s

   $ curl -s localhost:8080/metrics | grep pod_count
   # HELP pod_count Current number of pods
   # TYPE pod_count gauge
   pod_count{label_name="app",label_value="demo-app",phase="Running"} 2
   ```
1. When you're done, you can remove the cluster
    ```shell
    $ make kube-cleanup
    ```
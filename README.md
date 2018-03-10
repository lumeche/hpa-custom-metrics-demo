# Horizontal Pod Autoscaler based on custom metrics demo
This is a very simple and minimalistic demo of how the Horizontal pod Autoscaler works in Kubernetes. This demo is part of a presentation don in the kubernetes Montreal metup


## Requirements
1. Having a running k8s
2. Enable Custom metrics and Horizontal pod autoscaler
3. A running Prometheus server. In this example it's assumed Prometheus is running next to prometheus adapter. 

## Content

In each of the following two folders there are three files. The actual code of the server, a Dockerfile to build the image and a manifest file to deploy it in k8s
### scalable_server
Contains a simple server that has two points of entries. 
* The first one is a request that will hold the request for 10 seconds. That way we can simulate easily the load. 
* The second is returns the actual number of request being active. This is the one that the metric_server uses to build the metric
### metrics_server 
Contains a very simple metric server that gets the current number of request beeing served by the scalable_server. It then gets those metrics available to prometheus

##Service Monitor and HPA object
In the root folder there is a manifest.yaml file that contains the Service Monitor and the HPA objects. 
* The ServiceMonitor tells prometheus that he needs to grab the metrics from the custom_metrics server
* The HorizontalPodAutoscaler is the ones that tells k8s which component to scale and usign which metric to base from.

## DEMO

###Command to send traffic to the scalable server. 30 concurent connections

#### 10 channelsx
```bash
httperf --server 10.33.59.166  --port 28188 --uri /request --num-conns 30   --num-calls 999999 --rate 10
```
#### 14 channelsx
```bash
httperf --server 10.33.59.166  --port 28188 --uri /request --num-conns 30   --num-calls 999999 --rate 10
```
### To get the metrics publised by the metric_server
```bash
watch curl 'http://10.33.59.166:32416/metrics | grep utili'
```
### To get the list pod pods
```bash
watch  kubectl --context cia -n luist get pod
```

##TODO
Reduce the wait times after scaling up and down for the demo.
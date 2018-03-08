package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr           = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	scalableServer = flag.String("scalable-server", "http://localhost:8081", "The address to listen on for HTTP requests.")
	metric         = prometheus.NewGauge(prometheus.GaugeOpts{Name: "utilization", Help: "utilization for test server"})
	totalMap       map[string]float64
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(metric)
	totalMap = make(map[string]float64)
}

func getMetric() (string, float64) {
	url := *scalableServer + "/total"
	resp, err := http.Get(url)
	if err != nil {
		glog.Errorf("Error %v getting total", err)
	}
	defer resp.Body.Close()
	//TODO: Rrfactor this
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error("Error reading body")
	}
	hostnameAndTotal := strings.Split(string(bodyResp), ",")
	hostname := hostnameAndTotal[0]
	total, _ := strconv.ParseFloat(hostnameAndTotal[1], 64)
	return hostname, total

}

func calculateMaxTotal() float64 {
	var maxTotal float64
	for _, value := range totalMap {
		if maxTotal < value {
			maxTotal = value
		}
	}
	return maxTotal
}
func generateMetrics() {
	for true {
		hostname, total := getMetric()
		totalMap[hostname] = total
		maxTotal := calculateMaxTotal()
		metric.Set(float64(maxTotal))
		time.Sleep(5 * time.Second)
	}
}

func main() {
	flag.Parse()
	go generateMetrics()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

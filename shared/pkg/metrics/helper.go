package metrics

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusServer launches the promhttp handler on the given point
// This method is blocking.
func PrometheusServer(port int) {
	log.Printf("exposing prometheus metrics on /metrics on localhost:%d", port)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
	if err != nil {
		log.Fatalf("failed to serve prometheus http server: %v", err)
	}
}

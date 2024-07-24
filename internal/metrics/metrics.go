package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	messagesProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cm_enrich_messages_processed_total",
		Help: "Total number of messages processed",
	})
)

func init() {
	prometheus.MustRegister(messagesProcessed)
}

func SetupPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

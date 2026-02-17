package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	NotificationsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notifications_processed_total",
			Help: "Total processed notifications",
		},
	)

	NotificationsFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notifications_failed_total",
			Help: "Total failed notifications",
		},
	)
)

func Register() {
	prometheus.MustRegister(
		NotificationsProcessed,
		NotificationsFailed,
	)
}

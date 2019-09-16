package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type TrainTravel struct {
    to string
    from string
    depatureTime string
}

// Collector struct
type Collector struct {
	lateMetrics *prometheus.Desc
    cancelledMetrics *prometheus.Desc
    TrainTravel []TrainTravel
}

func newNSBCollector() *Collector {
	return &Collector{
		lateMetrics: prometheus.NewDesc("vy_arrival_after_schedule_seconds",
			"Vy train home is seconds late",
			[]string{"from", "to", "departure_scheduled"}, nil,
		),
		cancelledMetrics: prometheus.NewDesc("vy_cancelled_bool",
			"Vy train is cancelled",
			[]string{"from", "to", "departure_scheduled"}, nil,
		),
	}
}

// Describe metrics
func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.lateMetrics
	ch <- collector.cancelledMetrics
}

// Collect metrics
func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
    for _, travel := range collector.TrainTravel {
        vy := CallVy(travel)
	    ch <- prometheus.MustNewConstMetric(collector.lateMetrics, prometheus.CounterValue, IsVyLate(vy),
                travel.from, travel.to, travel.depatureTime)
	    ch <- prometheus.MustNewConstMetric(collector.cancelledMetrics, prometheus.CounterValue, IsCancelled(vy),
                travel.from, travel.to, travel.depatureTime)
    }
}

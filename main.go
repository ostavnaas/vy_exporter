package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//Create a new instance of the foocollector and
	//register it with the prometheus client.
	nsb := newNSBCollector()
    nsb.TrainTravel = []TrainTravel{
        TrainTravel{
            from: "Lilleby",
            to: "Hell",
            depatureTime: "15:10",
        },
        TrainTravel{
            from: "Hell",
            to: "Lilleby",
            depatureTime: "06:20",
        },
    }
	prometheus.MustRegister(nsb)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Vy exporter, use /metrics")

	})

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

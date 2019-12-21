#!/bin/bash

sudo docker network create prometheus
sudo docker run --network=prometheus --name vyexporter -d  vyexporter
sudo docker run -d --network=prometheus -p 9090:9090 -v $(pwd)/prometheus-data:/prometheus-data \
	prom/prometheus --config.file=/prometheus-data/prometheus.yml

sudo docker run -d --network=prometheus --name alertmanager -p 9093:9093 -v $(pwd)/prometheus-data:/alertmanager-data \
	prom/alertmanager

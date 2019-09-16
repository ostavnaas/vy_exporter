#!/bin/bash

sudo docker network create prometheus
sudo docker run --network=prometheus --name vyexporter -d  vyexporter
sudo docker run --network=prometheus -p 9090:9090 -v $(pwd)/prometheus-data:/prometheus-data \
	prom/prometheus --config.file=/prometheus-data/prometheus.yml

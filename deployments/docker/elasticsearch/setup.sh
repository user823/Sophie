#!/bin/bash

docker run -d --name sophie-es -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" -e ES_JAVA_OPTS="-Xmx512m -Xms256m -Xmn256m" elasticsearch:8.11.3

#kibana
docker run -d --name sophie-kibana -p 5601:5601 kibana:8.11.3
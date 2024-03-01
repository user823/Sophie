#!/bin/bash

checkVersion() {
    echo "Version = $1"
	echo $1 |grep -E "^[0-9]+\.[0-9]+\.[0-9]+" > /dev/null
    if [ $? = 0 ]; then
        return 1
    fi

	echo "Version $1 illegal, it should be X.X.X format(e.g. 4.5.0), please check released versions in 'https://archive.apache.org/dist/rocketmq/'"
    exit -1
}

helpMsg() {
    echo "Usage: setup version [-e]"
    echo "Options:"
    echo "  -e  setup rocketmq-dashboard"
    echo "  -p  set rocketmq-dashboard port, only when -e enabled"
    echo "  -H  set host ip"
    echo "  -q  quick start rocketmq in one container"
}

checkVersion $1
setup_dashboard
dashboard_port
host_ip

while getopts ":e" opt; do
  case $opt in
    e)
      setup_dashboard=true
    ;;
    p)
      if [-z "$OPTARG"]; then
        echo "-p option need a argument"
        exit 1
      fi
      dashboard_port=$OPTARG
    ;;
    h)
      if [-z "$OPTARG"]; then
        echo "-h option need a argument"
        exit 1
      fi
      host_ip=$OPTARG
    ;;
    q)
      echo "start sophie-rmq in one container"
      docker run -d --name sophie-rmq -p 9876:9876 -p 10911:10911 -p 8081:8081 apache/rocketmq:5.1.4 sophie-rmq
      exit 0
    :)
      # ignore other parameters
    ;;
  esac
done

# 开启namesrv 和 broker
docker run -d --name rmqnamesrv -p 9876:9876 apache/rocketmq:$1 bash mqnamesrv
docker run -d --name rmqbroker -p 10911:10911 apache/rocketmq:$1 bash mqbroker -n ${host_ip}:9876 --enable-proxy

# 开启rocketmq_dashboard
if [ "$setup_dashboard" = true ]; then
  docker run -d --name rmqdashboard -p ${dashboard_port}:${dashboard_port} -e "JAVA_OPTS=-Drocketmq.namesrv.addr=${host_ip}:9876" apache/rocketmq-dashboard
fi
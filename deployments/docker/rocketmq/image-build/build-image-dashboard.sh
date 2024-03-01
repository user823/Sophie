checkVersion() {
    echo "Version = $1"
	echo $1 |grep -E "^[0-9]+\.[0-9]+\.[0-9]+" > /dev/null
    if [ $? = 0 ]; then
        return 1
    fi

	echo "Version $1 illegal, it should be X.X.X format(e.g. 1.0.0), please check released versions in 'https://dist.apache.org/repos/dist/release/rocketmq/rocketmq-dashboard/'"
    exit -1
}


ROCKETMQ_DASHBOARD_VERSION=$1
dashboard_port=$2
if [-z "${dashboard_port}"]; then
  dashboard_port=8080
fi

checkVersion $ROCKETMQ_DASHBOARD_VERSION

# Build rocketmq
# 通过dashboard_port 指定开放开放的端口
docker build --no-cache -f dockerfile-dashboard -t apache/rocketmq-dashboard:latest --build-arg version=${ROCKETMQ_DASHBOARD_VERSION} --build-arg dashboard_port=${dashboard_port} .
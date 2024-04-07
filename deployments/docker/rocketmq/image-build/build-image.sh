#!/bin/bash

checkVersion() {
    echo "Version = $1"
	echo $1 |grep -E "^[0-9]+\.[0-9]+\.[0-9]+" > /dev/null
    if [ $? = 0 ]; then
        return 1
    fi

	echo "Version $1 illegal, it should be X.X.X format(e.g. 5.1.4), please check released versions in 'https://archive.apache.org/dist/rocketmq/'"
    exit -1
}

if [ $# -ne 1 ]; then
    echo -e "Usage: sh $0 Version(format e.g. 5.1.4)"
    exit -1
fi

ROCKETMQ_VERSION=$1
BASE_IMAGE=$2

checkVersion $ROCKETMQ_VERSION

# Build rocketmq
docker build --no-cache -f dockerfile -t sophie/rocketmq:${ROCKETMQ_VERSION} --build-arg version=${ROCKETMQ_VERSION} .
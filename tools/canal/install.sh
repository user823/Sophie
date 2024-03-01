#!/bin/bash

canal-server-url=`https://github.com/alibaba/canal/releases/download/canal-1.1.7/canal.deployer-1.1.7.tar.gz`
canal-adapter-url=`https://github.com/alibaba/canal/releases/download/canal-1.1.7/canal.adapter-1.1.7.tar.gz`
canal-admin-url=`https://github.com/alibaba/canal/releases/download/canal-1.1.7/canal.admin-1.1.7.tar.gz`

wget canal-server-url
mkdir canal-server
tar -zxvf canal.deployer-1.1.7.tar.gz -C canal-server && rm canal.deployer-1.1.7.tar.gz

wget canal-adapter-url
mkdir canal-adapter
tar -zxvf canal.adapter-1.1.7.tar.gz -C canal-adapter && rm canal.adapter-1.1.7.tar.gz

wget canal-admin-url
mkdir canal-admin
tar -zxvf canal.admin-1.1.7.tar.gz -C canal-admin && rm canal.admin-1.1.7.tar.gz
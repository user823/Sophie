#!/bin/bash

canal_server_url="https://github.com/alibaba/canal/releases/download/canal-1.1.7/canal.deployer-1.1.7.tar.gz"
canal_adapter_url="https://github.com/alibaba/canal/releases/download/canal-1.1.7/canal.adapter-1.1.7.tar.gz"
canal_admin_url="https://github.com/alibaba/canal/releases/download/canal-1.1.7/canal.admin-1.1.7.tar.gz"

wget $canal_server_url
mkdir canal-server
tar -zxvf canal.deployer-1.1.7.tar.gz -C canal-server && rm canal.deployer-1.1.7.tar.gz

wget $canal_adapter_url
mkdir canal-adapter
tar -zxvf canal.adapter-1.1.7.tar.gz -C canal-adapter && rm canal.adapter-1.1.7.tar.gz

wget $canal_admin_url
mkdir canal-admin
tar -zxvf canal.admin-1.1.7.tar.gz -C canal-admin && rm canal.admin-1.1.7.tar.gz

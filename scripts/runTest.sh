#!/bin/bash

#运行基本环境：etcd (ps.rabbitmq已在安装时设置为开机自启动）
etcd &


cd ../cmd/apiserver
go build
./apiserver
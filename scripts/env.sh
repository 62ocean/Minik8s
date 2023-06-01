#!/bin/bash
#更换apt源
sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
apt-get update
apt-get install -y wget
echo "—————————————————成功更换apt源—————————————————"

#安装go
wget -c https://dl.google.com/go/go1.20.4.linux-amd64.tar.gz -O - | tar -xz -C /usr/local
echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile
source ~/.profile
#验证
export PATH=$PATH:/usr/local/go/bin
go version

# 修改镜像源
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
echo "—————————————————完成go语言安装—————————————————"

#安装etcd
wget https://storage.googleapis.com/etcd/v3.4.26/etcd-v3.4.26-linux-amd64.tar.gz
tar zxvf etcd-v3.4.26-linux-amd64.tar.gz
mv etcd-v3.4.26-linux-amd64 etcd
cd etcd
cp etcd /usr/local/bin
cp etcdctl /usr/local/bin
cd ..
cp -r etcd /usr/local/etcd
rm -rf etcd
echo "—————————————————完成etcd安装—————————————————"

#安装rabbitmq
#安装erlang语言
apt-get install -y erlang-nox
#添加公钥
wget -O- https://www.rabbitmq.com/rabbitmq-release-signing-key.asc | apt-key add -
apt-get update
#安装rabbitmq
apt-get install -y rabbitmq-server systemd
#设置开机自启动
systemctl enable rabbitmq-server
#查看rabbitmq状态（此时应该正在运行了）
systemctl status rabbitmq-server
echo "———————————————完成rabbitmq安装—————————————————"

# 安装nginx
apt-get install -y nginx
echo "———————————————完成nginx安装—————————————————"


# 关闭原本的dns，此时使用域名访问网络将失效，之后就可以开始跑minik8s啦！(在构建脚本里就不跑这个了，测试的时候再关)
#systemctl stop systemd-resolved
#systemctl disable systemd-resolved

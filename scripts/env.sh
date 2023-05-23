#!/bin/bash
#安装go
wget -c https://dl.google.com/go/go1.20.4.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
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
wget -O- https://www.rabbitmq.com/rabbitmq-release-signing-key.asc | sudo apt-key add -
apt-get update
#安装rabbitmq
apt-get install -y rabbitmq-server
#设置开机自启动
systemctl enable rabbitmq-server
#查看rabbitmq状态（此时应该正在运行了）
systemctl status rabbitmq-server
echo "———————————————完成rabbitmq安装—————————————————"

# 安装nginx
apt-get install -y nginx
echo "———————————————完成nginx安装—————————————————"


# 关闭原本的dns，此时使用域名访问网络将失效，之后就可以开始跑minik8s啦！
systemctl stop systemd-resolved
systemctl disable systemd-resolved

exit_script
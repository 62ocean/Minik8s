# kubectl命令

> 使用：在cmd文件夹下，执行`go build -o kubectl`命令，得到`kubectl`可执行文件。
目前所有的命令都没有调用实际接口，只是输出了命令信息。

## help
```shell
./kubectl --help
```

## create
### 创建pod
```shell
./kubectl create -f demo-pod.yaml
```

## delete
### 使用demo-pod.yaml中指定的资源类型和名称删除pod
```shell
./kubectl delete -f demo-pod.yaml
```

## describe
### 获得pod的运行状态
```shell
./kubectl describe pod demo-pod 
```
### 获得service的运行状态
```shell
./kubectl describe service demo-service
```
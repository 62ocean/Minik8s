# ReplicaSet实现文档

> Q1: apiserver启动后会将所有pod/service信息保存在内存中吗？还是每次需要这些信息的时候调函数从etcd中拿呢？需要缓存吗？  
> Q2: 对etcd的操作需要上锁吗？  
> Q3: 新建pod和删除pod时要先更改etcd再进行操作吗（即以etcd中状态为准）？

## 实现方式

1. replicaset controller被动监听pods变化，如发生改变则启动一个工作线程进行处理（k8s的实现方式，类似中断）。
2. replicaset controller每隔一段时间主动检查一遍所有replicaset的期望状态是否与实际状态一致，同时进行处理（类似轮询）。

方式1性能更优，方式2实现简单。目前采用第2种实现方式，后续可以改为第1种。

## 调用层级

`server` ---(启动时创建)---> `replicaset controller`  
`replicaset controller` ---(监听)---> `apiserver` ---> `pods status`   
- 实际状态与期望状态一致，不做任何处理   
- 实际状态与期望状态不同 ---(调接口)---> `apiserver` ---> `create/delete pods`   


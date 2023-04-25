### ETCD存储各个API对象的格式
> 依据K8S原本实现的格式进行存储

所有资源的组织格式都是：
```
/registry/{resource_name}/{namespace}/{resource_instance}
```
各个资源的id使用google/uuid包来生成，记得提前运行如下命令
```bash
go get github.com/google/uuid
```

#### 集群数据
##### node

```
/registry/node/<node-ip-1>
/registry/node/<node-ip-2>
/registry/node/<node-ip-3>
```

#### k8s对象数据

##### namespace

```bash
/registry/namespaces/default
/registry/namespaces/game
/registry/namespaces/kube-node-lease
/registry/namespaces/kube-public
/registry/namespaces/kube-system
```



##### k8s常见对象：

```bash
/registry/{resource}/{namespace}/{resource_name}
```

```bash
# deployment
/registry/deployments/default/game-2048
/registry/deployments/kube-system/prometheus-operator

# replicasets
/registry/replicasets/default/game-2048-c7d589ccf

# pod
/registry/pods/default/game-2048-c7d589ccf-8lsbw

# statefulsets
/registry/statefulsets/kube-system/prometheus-k8s

# daemonsets
/registry/daemonsets/kube-system/kube-proxy

# secrets
/registry/secrets/default/default-token-tbfmb

# serviceaccounts
/registry/serviceaccounts/default/default
```



##### service

```bash
# service
/registry/services/specs/default/game-2048

# endpoints
/registry/services/endpoints/default/game-2048
```

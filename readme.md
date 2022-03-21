### Operator创建过程
```shell
mkdir test-operator && cd test-operator

go mod init test-operator

operator-sdk init

operator-sdk create api --group ship --version v1beta1 --kind Frigate
```


### Operator启动过程

#### 1.创建crd
```shell
make manifests  # 会创建config/crd下的default文件

make generate

make install  # or  kubectl apply -k config/crd/
```

> 可以使用 `kubectl apply -k config/crd` 生成crd文件，但是当这样使用时
> 如果你在type文件里修改了结构体字段，需要在config/crd/bases文件里手动添加这些字段


#### 2.创建rbac
```shell
kubectl apply -k config/rbac  # 要注意提前修复对应文件里的ns和name

```


#### 3.build镜像
```shell
```


#### 4.创建controller manager


#### 5.创建sample operator
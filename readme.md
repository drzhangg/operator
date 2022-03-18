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


#### 2.创建rbac
```shell
kubectl apply -k config/rbac  # 要注意提前修复对应文件里的ns和name

```


#### 3.build镜像
```shell
```


#### 4.创建controller manager


#### 5.创建sample operator
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
make matis

make install
```


#### 2.创建rbac


#### 3.build镜像


#### 4.创建controller manager


#### 5.创建sample operator
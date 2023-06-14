# image-syncer

此`image-syncer` 是一个定时docker镜像同步工具，可用来进行多对多的镜像仓库同步，支持目前绝大多数主流的docker镜像仓库服务，接收同步hook信息，定时同步

[English](./README.md) | 简体中文

## Features

- 支持多对多镜像仓库同步
- 支持基于Docker Registry V2搭建的docker镜像仓库服务 (如 Docker Hub、 Quay、 阿里云镜像服务ACR、 Harbor等)
- 同步只经过内存和网络，不依赖磁盘存储，同步速度快
- 增量同步, 通过对同步过的镜像blob信息落盘，不重复同步已同步的镜像
- 并发同步，可以通过配置文件调整并发数
- 自动重试失败的同步任务，可以解决大部分镜像同步中的网络抖动问题
- 不依赖docker以及其他程序
- 通过docker hook工具将同步镜像推送至此，此工具会记录信息进而定时同步

## 使用



### 手动编译

```
go get github.com/zxzixuanwang/image-syncer
cd $GOPATH/github.com/zxzixuanwang/image-syncer

# This will create a binary file named image-syncer
make
```
### 参数
```bash
# example 
curl  -i http://localhost:8080/images/sync/hook\?name\=reponame/namespace/imagename\&tag\=1.0.3 -u $username:$password

```

### 使用用例

```shell
# 默认启动，会读取configs文件夹下，sync.yaml文件
./image-syncer 

# 指定配置文件启动
./image-syncer -c configs/sync.yaml
```

<!-- 
### 同步镜像到ACR

ACR(Ali Container Registry) 是阿里云提供的容器镜像服务，ACR企业版(EE)提供了企业级的容器镜像、Helm Chart 安全托管能力，推荐安全需求高、业务多地域部署、拥有大规模集群节点的企业级客户使用。

这里会将quay.io上的一些镜像同步到ACR企业版，作为使用用例。

#### 创建企业版ACR

1. [创建容器镜像服务]()
2.  -->

### 配置文件

在 v1.2.0 版本之后，image-syncer 的配置文件支持JSON和YAML两种格式，并且支持将原config文件替换为一个认证信息文件和一个镜像同步文件。详细的配置文件示例可在目录 [example](./example) 下找到，旧版本的配置文件格式（auth 和 images 字段放在一起的版本，通过 --config 参数指定）也是兼容的，目录下 `config.json` 为示例。

#### 认证信息

`auth.json` 包含了所有仓库的认证信息

```java
{  
    // 认证字段，其中每个对象为一个registry的一个账号和
    // 密码；通常，同步源需要具有pull以及访问tags权限，
    // 同步目标需要拥有push以及创建仓库权限，如果没有提供，则默认匿名访问
    
    "quay.io": {    // 支持 "registry" 和 "registry/namespace"（v1.0.3之后的版本） 的形式，需要跟下面images中的registry(registry/namespace)对应
                    // images中被匹配到的的url会使用对应账号密码进行镜像同步, 优先匹配 "registry/namespace" 的形式
        "username": "xxx",               // 用户名，可选，（v1.3.1 之后支持）valuse 使用 "${env}" 或者 "$env" 类型的字符串可以引用环境变量
        "password": "xxxxxxxxx",         // 密码，可选，（v1.3.1 之后支持）valuse 使用 "${env}" 或者 "$env" 类型的字符串可以引用环境变量
        "insecure": true                 // registry是否是http服务，如果是，insecure 字段需要为true，默认是false，可选，支持这个选项需要image-syncer版本 > v1.0.1
    },
    "registry.cn-beijing.aliyuncs.com": {
        "username": "xxx",
        "password": "xxxxxxxxx"
    },
    "registry.hub.docker.com": {
        "username": "xxx",
        "password": "xxxxxxxxxx"
    },
    "quay.io/coreos": {                       
        "username": "abc",              
        "password": "xxxxxxxxx",
        "insecure": true  
    }
}
```

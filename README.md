# image-syncer

此`image-syncer` 是一个经过修改过定时docker镜像同步工具，可用来进行多对多的镜像仓库同步，接收同步hook信息，定时同步。

### 手动编译

```
go get github.com/zxzixuanwang/image-syncer
cd $GOPATH/github.com/zxzixuanwang/image-syncer

# This will create a binary file named image-syncer
make
```
### HTTP参数
```bash
# example 
curl  -i http://localhost:8080/images/sync/hook\?name\=reponame/namespace/imagename\&tag\=1.0.3 -u $username:$password

```
### rabbitmq 消息
```golang
// struct
type RabbitMqData struct {
	Name string
	Tag  string
}
```


### 命令用例

```shell
# 默认启动，会读取configs文件夹下，sync.yaml文件
./image-syncer 

# 指定配置文件启动
./image-syncer -c configs/sync.yaml

# auth.yaml
文件需要存放到运行目录下面
```
#### 认证信息

认证信息中可以同时描述多个 registry（或者 registry/namespace）对象，一个对象可以包含账号和密码，其中，密码可能是一个 TOKEN

> 注意，通常镜像源仓库需要具有 pull 以及访问 tags 权限，镜像目标仓库需要拥有 push 以及创建仓库权限；如果对应仓库没有提供认证信息，则默认匿名访问

认证信息文件通过 `--auth` 参数传入，具体文件样例可以参考 [auth.yaml](examples/auth.yaml) 和 [auth.json](examples/auth.json)，这里以 [auth.yaml](examples/auth.yaml) 为例：

```yaml
quay.io: 
  username: xxx
  password: xxxxxxxxx
  insecure: true # 可选，（v1.0.1 之后支持）registry是否是http服务，如果是，insecure 字段需要为 true，默认是 false
registry.cn-beijing.aliyuncs.com:
  username: xxx # 可选，（v1.3.1 之后支持）value 使用 "${env}" 或者 "$env" 形式可以引用环境变量
  password: xxxxxxxxx # 可选，（v1.3.1 之后支持）value 使用 "${env}" 或者 "$env" 类型的字符串可以引用环境变量
quay.io/coreos:
  username: abc
  password: xxxxxxxxx
  insecure: true
```
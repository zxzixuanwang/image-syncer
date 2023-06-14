# image-syncer


This `image-syncer` is a timed Docker image synchronization tool that can be used for multi to many image warehouse synchronization. With `image-syncer` you can synchronize docker images from some source registries to target registries, which include most popular public docker registry services.

English | [简体中文](./README-zh_CN.md)

## Features

- Support for many-to-many registry synchronization
- Supports docker registry services based on Docker Registry V2 (e.g., Alibaba Cloud Container Registry Service, Docker Hub, Quay.io, Harbor, etc.)
- Network & Memory Only, don't rely on large disk storage, fast synchronization
- Incremental Synchronization, use a disk file to record the synchronized image blobs' information
- Concurrent Synchronization, adjustable goroutine numbers
- Automatic Retries of Failed Sync Tasks, to resolve the network problems while synchronizing
- Doesn't rely on Docker daemon or other programs
- Push the synchronized image here through the Docker Hook tool, whic

## Usage


### Compile Manually

```
go get github.com/zxzixuanwang/image-syncer
cd $GOPATH/github.com/zxzixuanwang/image-syncer

# This will create a binary file named image-syncer
make
```

### Parameters
```bash
# example 
curl  -i http://localhost:8080/images/sync/hook\?name\=reponame/namespace/imagename\&tag\=1.0.3 -u $username:$password

```

### Example

```shell
# By default, it will read the sync.yaml file in the configs folder
./image-syncer 

# Specify profile startup
./image-syncer -c configs/sync.yaml

### Configure Files

After v1.2.0, image-syncer supports both YAML and JSON format, and origin config file can be split into "auth" and "images" file. A full list of examples can be found under [example](./example), meanwhile the older version of configuration file is still supported via --config flag.

#### Authentication file

Authentication file holds all the authentication information for each registry, the following is an example of `auth.json`

```java
{               
    // Authentication fields, each object has a URL as key and a username/password pair as value, 
    // if authentication object is not provided for a registry, access to the registry will be anonymous.
        
    "quay.io": {        // This "registry" or "registry/namespace" string should be the same as registry or registry/namespace used below in "images" field.  
                            // The format of "registry/namespace" will be more prior matched than "registry"
        "username": "xxx",       // Optional, if the value is a string of "${env}" or "$env", image-syncer will try to find the value in environment variables, after v1.3.1       
        "password": "xxxxxxxxx", // Optional, if the value is a string of "${env}" or "$env", image-syncer will try to find the value in environment variables, after v1.3.1
        "insecure": true         // "insecure" field needs to be true if this registry is a http service, default value is false, version of image-syncer need to be later than v1.0.1 to support this field
    },
    "registry.cn-beijing.aliyuncs.com": {
        "username": "xxx",
        "password": "xxxxxxxxx"
    },
    "registry.hub.docker.com": {
        "username": "xxx",
        "password": "xxxxxxxxxx"
    },
    "quay.io/coreos": {     // "registry/namespace" format is supported after v1.0.3 of image-syncer     
        "username": "abc",              
        "password": "xxxxxxxxx",
        "insecure": true  
    }
}
```



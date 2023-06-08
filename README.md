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

#### Image sync configuration file

Image sync configuration file defines all the image synchronization rules, the following is an example of `images.json`

```java
{
    // Rules of image synchronization, each rule is a kv pair of source(key) and destination(value). 

    // The source of each rule should not be empty string.

    // If you need to synchronize images from one source to multi destinations, add more rules.

    // Both source and destination are docker image url (registry/namespace/repository:tag), 
    // with or without tags.

    // For both source and destination, if destination is not an empty string, "registry/namespace/repository" 
    // is needed at least.
    
    // You cannot synchronize a whole namespace or a registry but a repository for one rule at most.

    // The repository name and tag of destination can be deferent from source, which works like 
    // "docker pull + docker tag + docker push"

    "quay.io/coreos/kube-rbac-proxy": "quay.io/ruohe/kube-rbac-proxy",
    "xxxx":"xxxxx",
    "xxx/xxx/xx:tag1,tag2,tag3":"xxx/xxx/xx"

    // If a source doesn't include tags, it means all the tags of this repository need to be synchronized,
    // destination should not include tags at this moment.
    
    // Each source can include more than one tags, which is split by comma (e.g., "a/b/c:1", "a/b/c:1,2,3").

    // If a source includes just one tag (e.g., "a/b/c:1"), it means only one tag need to be synchronized;
    // at this moment, if the destination doesn't include a tag, synchronized image will keep the same tag.
    
    // When a source includes more than one tag (e.g., "a/b/c:1,2,3"), at this moment,
    // the destination should not include tags, synchronized images will keep the original tags.
    // e.g., "a/b/c:1,2,3":"x/y/z".
    
    // When a destination is an empty string, source will be synchronized to "default-registry/default-namespace"
    // with the same repository name and tags, default-registry and default-namespace can be set by both parameters
    // and environment variable.
}	
```



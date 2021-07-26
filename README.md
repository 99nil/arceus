# Pokémon - arceus(读音 阿尔宙斯)

![LICENSE](https://img.shields.io/github/license/zc2638/arceus.svg?style=flat-square&color=blue)
[![Go Reference](https://pkg.go.dev/badge/github.com/zc2638/arceus.svg)](https://pkg.go.dev/github.com/zc2638/arceus)
[![Go Report Card](https://goreportcard.com/badge/github.com/zc2638/arceus)](https://goreportcard.com/report/github.com/zc2638/arceus)
![Main CI](https://github.com/zc2638/arceus/workflows/Main%20CI/badge.svg)
<a target="_blank" href="https://qm.qq.com/cgi-bin/qm/qr?k=d_FApC9aD6o6XZ2LR0zx5uO5Z642bP6M&jump_from=webapi"><img border="0" src="https://pub.idqqimg.com/wpa/images/group.png" alt="99nil" title="99nil"></a>

## 定义
  可结构化内容构造器
  
## 依赖
Go Version 1.16+

## 用途
  - 可用Devops中，快速构造各类Yaml文件 （例k8s yaml、java 配置yaml ...），提升研发工作效率
  - 可用于作为业务组件编排配置生成
  - etc...

## 使用
[使用手册](https://github.com/99nil/arceus/blob/main/docs/help.md)

## TODO 

- 接口模式优化
- 格式转换(e.g. json <=> yaml)

## Run
### Local
```shell
go run github.com/zc2638/arceus/cmd
```

### Docker
基础启动
```shell
docker run --name arceus -d -p 2638:2638 zc2638/arceus:latest
```
挂载启动
```shell
docker run --name arceus -d -p 2638:2638 -v ~/docker/arceus:/etc/arceus zc2638/arceus:latest
```
使用镜像执行QuickStart
```shell
docker run --rm -it \
 -v ~/docker/arceus:/etc/arceus \
 -v ~/docker/arceus/examples:/work/examples \
 zc2638/arceus:latest \
 sh -c './arceus apply -f /work/examples/template/nginx.yaml \
 && ./arceus apply -f /work/examples/quickstart/app/app-rule.yaml \
 && ./arceus qs -f /work/examples/quickstart/app/app.yaml -o /etc/arceus/output'
```

## Build
### Build/Update UI
```shell
make ui
```

### Build image
```shell
make docker
```

# Pokémon - arceus(读音 阿尔宙斯)

![LICENSE](https://img.shields.io/github/license/zc2638/arceus.svg?style=flat-square&color=blue)
[![Go Reference](https://pkg.go.dev/badge/github.com/zc2638/arceus.svg)](https://pkg.go.dev/github.com/zc2638/arceus)
[![Go Report Card](https://goreportcard.com/badge/github.com/zc2638/arceus)](https://goreportcard.com/report/github.com/zc2638/arceus)
![Main CI](https://github.com/zc2638/arceus/workflows/Main%20CI/badge.svg)


# 定义
  可视化结构数据构造器

# 用途
  可用于快速构造各类Yaml文件 （例k8s yaml、java 配置yaml ...）
  可用于作为业务组件编排配置构造器

# 
## TODO 

- 命令行模式
- 接口模式优化
- 格式转换(e.g. json <=> yaml)

## Run
### Local
```shell
go run github.com/zc2638/arceus/cmd
```

### Docker
```shell
docker run --name arceus -d -p 2638:2638 zc2638/arceus:latest
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
# arceus

![LICENSE](https://img.shields.io/github/license/zc2638/arceus.svg?style=flat-square&color=blue)
[![Go Reference](https://pkg.go.dev/badge/github.com/zc2638/arceus.svg)](https://pkg.go.dev/github.com/zc2638/arceus)
[![Go Report Card](https://goreportcard.com/badge/github.com/zc2638/arceus)](https://goreportcard.com/report/github.com/zc2638/arceus)
![Main CI](https://github.com/zc2638/arceus/workflows/Main%20CI/badge.svg)

Pokémon - arceus

Structured configuration generator.

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
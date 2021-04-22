# arceus

Pokémon - arceus

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
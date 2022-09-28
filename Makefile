.PHONY: ui
ui:
	./build/ui.sh

.PHONY: docker
docker:
	docker build -t zc2638/arceus -f build/Dockerfile .

.PHONY: test
test:
	golangci-lint run ./... && go test ./...

.PHONY: build
build:
	./build/build.sh

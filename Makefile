.PHONY: ui
ui:
	./build/update.sh

.PHONY: docker
docker:
	docker build -t zc2638/arceus -f build/Dockerfile .
OBJECTS=server.out

IMG=golang-boilerplate-be
IMG_TAG=v1

CONTAINER_REGISTRY = docker.io
USER = duyle95

staticcheck:
	staticcheck ./...

build:
	swag init -g ./cmd/service/main.go
	go build -o $(OBJECTS) cmd/service/main.go

run:
	make build
	./$(OBJECTS)

docker-build:
	docker build -t $(IMG):$(IMG_TAG) .

docker-run:
	cd deployments && docker compose up

docker-push:
	docker tag $(IMG):$(IMG_TAG) $(USER)/$(IMG):$(IMG_TAG)
	docker push $(CONTAINER_REGISTRY)/$(USER)/$(IMG):$(IMG_TAG)

clean:
	rm $(OBJECTS)

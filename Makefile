# version = $(shell git describe --tags)
version = "v0.1.0"
# ldflags = -ldflags "-s -w -X main.version=${version}"

build:
	go build ${ldflags} -v ./cmd/rkcd
	go build ${ldflags} -v ./cmd/rkcd-cli

run:
	./rkcd

test:
	go test -race ./...

dcu:
	docker-compose up -d


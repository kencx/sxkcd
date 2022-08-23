# version = $(shell git describe --tags)
version = "v0.1.0"
ldflags = -ldflags "-s -w -X main.version=${version} -X github.com/kencx/rkcd/http.version=${version}"

.PHONY: build run test dcu clean

build:
	go build ${ldflags} -v .

run:
	./rkcd

test:
	go test -race ./...

dcu:
	docker-compose up -d

dbuild:
	docker-compose up -d --build

clean:
	rm -rf rkcd

version = $(shell git describe --tags)
ldflags = -ldflags "-s -w -X main.version=${version}"

.PHONY: build dbuild test dcu clean

build:
	go build ${ldflags} -v .

dbuild: Dockerfile
	docker build . -t ghcr.io/kencx/sxkcd:${version}

test:
	go test -race ./data

dcu: docker-compose.yml
	docker-compose up -d

clean:
	rm -rf sxkcd
	rm -rf ui/build/*

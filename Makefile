build:
	go build -v .

run:
	./rkcd

test:
	go test -race ./...

dcu:
	docker-compose up -d

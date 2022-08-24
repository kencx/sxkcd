# version = $(shell git describe --tags)
version = "v0.1.0"
ldflags = -ldflags "-s -w -X main.version=${version}"

.PHONY: build dbuild run test dcu clean deploy destroy

build:
	go build ${ldflags} -v .

dbuild: docker-compose.yml
	docker-compose up -d --build

run:
	./rkcd

test:
	go test -race ./data

dcu: docker-compose.yml
	docker-compose up -d

deploy: deploy/terraform.tfstate
	cd deploy && terraform apply

destroy:
	cd deploy && terraform destroy

clean:
	rm -rf rkcd

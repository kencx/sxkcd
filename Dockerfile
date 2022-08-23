FROM node:18-alpine3.15 as frontend
WORKDIR /ui
COPY ui/package.json ui/package-lock.json ./
RUN npm ci --quiet
COPY ui/ .
RUN npm run build

FROM golang:1.18-alpine3.15 as builder
WORKDIR /app
ENV CGO_ENABLED=0 GOFLAGS="-ldflags=-s -w"
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend /ui/build ./ui/build

RUN go vet -v
RUN go build -v .

FROM alpine:3.15
# requires buildkit: DOCKER_BUILDKIT=1
COPY --from=builder --chmod=+x /app/rkcd /app/entrypoint.sh ./app/data/comics.json ./
EXPOSE 6380
ENTRYPOINT [ "/entrypoint.sh" ]

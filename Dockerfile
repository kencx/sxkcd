FROM node:18-alpine3.18 as frontend
WORKDIR /ui
COPY ui/package.json ui/package-lock.json ./
RUN npm ci --quiet
COPY ui/ .
RUN npm run build

FROM golang:1.20-alpine3.18 as builder
WORKDIR /app
ENV CGO_ENABLED=0 GOFLAGS="-ldflags=-s -w"
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend /ui/build ./ui/build

RUN go vet -v && go build -v .

FROM alpine:3.18
LABEL maintainer="kencx"

WORKDIR /
# requires buildkit: DOCKER_BUILDKIT=1
COPY --from=builder --chmod=+x /app/sxkcd /app/entrypoint.sh ./
EXPOSE 6380
ENTRYPOINT [ "/entrypoint.sh" ]

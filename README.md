# rkcd

Yet another [XKCD](https://xkcd.com) search engine.

<!-- ## Querying -->

## Usage
```bash
usage: rkcd [server|download] [OPTIONS] [FILE]

  Options:
    -v, --version   Version info
    -h, --help	    Show help

  server:
    -f, --file      Read data from file
    -p, --port      Server port
    -r, --redis     Redis connection URI [host:port]

  download:
    -l, --latest    Get latest comic number
    -n, --num       Download single comic by number
    -f, --file	    Download all comics to file
```

rkcd requires a dataset to start the server. Download the full set of xkcd comics

```bash
$ rkcd download -f data/comics.json
```

To start your own instance of rkcd, use the provided `docker-compose.yml` which starts two containers: Redis and rkcd.

```bash
$ docker-compose up -d
```

Alternatively, rkcd can be run directly on the host with a local instance of Redis and the `rkcd` binary which must be built from source.

```bash
$ redis-server
$ rkcd server -p 6380 -r localhost:6379 -f data/comics.json
```

## Development
rkcd is built with

- Go 1.17
- [Redis](https://redis.io/) w/[RediSearch](https://redis.io/docs/stack/search/) and [RedisJSON](https://redis.io/docs/stack/json/)
- [Svelte](https://svelte.dev/) 3.46.0 and [Sveltekit](https://kit.svelte.dev/) 1.0.0-next.405
- [picocss](https://picocss.com/) v1.5.3
- Node.js v18
- Docker (optional)

#### Build from Source
Because the frontend's static files are embedded in the Go binary, we must generate them prior to building the binary.

```bash
$ git clone https://github.com/kencx/rkcd.git
$ cd ui && npm ci --quiet

# Build static files in ui/build
$ cd ui && npm run build

# Build binary
$ go mod download
$ make build
```

#### Local Development
```bash
# Start frontend
$ cd ui
$ npm install
$ npm run dev

# Start redis docker container
$ docker-compose up -d redis

# Start rkcd server
$ go mod download
$ make build
$ ./rkcd server --port 6380 --redis localhost:6379 --file data/comics.json
```

Visit `localhost:5173`.

<!-- ## Deployment -->

## License

[MIT](LICENSE)

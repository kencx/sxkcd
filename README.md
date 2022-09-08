# sxkcd
Search for that [xkcd](https://xkcd.com) you swear you remember...

## Introduction
sxkcd is built with Go, Redis and Svelte.

Data is collected from xkcd and [explainxkcd](https://explainxkcd.com), and indexed with
Redis. Redis was chosen for its simplicity and high performance. Paired with RediSearch
and RediJSON, it provides full-text search capabilities with an extensive query syntax.
With a relatively small dataset, it processes queries really quickly in <40ms.

The frontend is a reactive static site built with Svelte. It fetches responses to your
queries from the backend Go server, which provides a query API for the in-memory database.

Try it out [here](https://sxkcd.cheo.dev)!

>To update the database, a cronjob runs `scripts/refresh.sh` every Tuesday, Thursday and
>Saturday morning to download and reindex the latest dataset. So expect some delay for
>the newest comic.

## Usage
```bash
usage: sxkcd [server|download] [OPTIONS] [FILE]

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

sxkcd requires a dataset to start the server. Download the full set of xkcd comics with

```bash
$ sxkcd download -f data/comics.json
```

To start your own instance of sxkcd, use the provided `docker-compose.yml` which starts
two Docker containers: Redis and sxkcd.

```bash
$ docker-compose up -d
```

Alternatively, sxkcd can be run directly on the host with a local instance of Redis and the `sxkcd` binary, which must be built from source.

```bash
$ redis-server
$ sxkcd server -p 6380 -r localhost:6379 -f data/comics.json
```

## Querying
sxkcd supports union, negation, prefix matching and filtering by custom date ranges
```text
# foo OR bar
> foo|bar

# exclude foo
> -foo

# prefix matching
> foo*

# filter by comic number
> #420

# comic number range
> #420-425

# filter by dates in ISO-8601 format
> @date: 2022-01-01, 2022-02-01

# from date to present
> @date: 2022-08-01
```

## Development
sxkcd is built with

- Go 1.17
- [Redis](https://redis.io/) w/[RediSearch](https://redis.io/docs/stack/search/) and [RedisJSON](https://redis.io/docs/stack/json/)
- [Svelte](https://svelte.dev/) 3.46.0 and [Sveltekit](https://kit.svelte.dev/) 1.0.0-next.405
- [picocss](https://picocss.com/) v1.5.3
- Node.js v18
- Docker (optional)

#### Build from Source
Because the frontend's static files are embedded in the Go binary, we must generate them prior to building:

```bash
$ git clone https://github.com/kencx/sxkcd.git
$ cd ui && npm ci --quiet

# Build static files in ui/build
$ cd ui && npm run build

# Build binary
$ go mod download
$ make build
```

#### Build Docker Image
You can also choose to build a local docker image after generating the static files.

```bash
$ make dbuild
```

#### Local Development
```bash
# Start frontend
$ cd ui
$ npm install
$ npm run dev

# Start redis docker container
$ docker-compose up -d redis

# Start sxkcd server
$ go mod download
$ make build
$ ./sxkcd server --port 6380 --redis localhost:6379 --file data/comics.json
```

Visit `localhost:5173`.

## Acknowledgements
>sxkcd is entirely inspired by [classes.wtf](https://github.com/ekzhang/classes.wtf)

[MIT](LICENSE)

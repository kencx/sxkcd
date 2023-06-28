# sxkcd
>Search for that [xkcd](https://xkcd.com) you swear you remember...

sxkcd is a simple webserver powered by Go, Redis and
[Svelte](https://svelte.dev).

All data is scraped from xkcd and [explainxkcd](https://explainxkcd.com), and
indexed with Redis. Paired with [RediSearch](https://redis.io/docs/stack/search)
and [RediJSON](https://redis.io/docs/stack/search/), Redis provides full-text
search and an extensive query syntax to the indexed data. Every day, sxkcd
automatically updates the database with the newest comic, if it doesn't already
exist.

Try it out [here](https://sxkcd.cheo.dev)!

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
    -i, --reindex   Reindex existing data with new file

  download:
    -l, --latest    Get latest comic number
    -n, --num       Download single comic by number
    -f, --file	    Download all comics to file
```

Start your own instance of sxkcd (and Redis) with the provided `docker-compose.yml`:

```console
$ docker-compose -d
```

### Local Setup

To run locally:

```console
# download all comics
$ sxkcd download -f data/comics.json

# start a redis instance with data persistence
$ redis-server --appendonly yes
$ sxkcd server -p 6380 -r localhost:6379 -f data/comics.json
```

If Redis is started with persistence, sxkcd can be restarted without any data
files. If we wish to reindex all data in the database with a new file, we can
run sxkcd with the `--reindex` flag:

```console
$ sxkcd server -p 6380 -r localhost:6379 -f data/new.json --reindex
```

This will replace all existing data with data in the new file.

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

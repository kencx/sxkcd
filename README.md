# sxkcd

S(earch)xkcd is a [xkcd](https://xkcd.com) search engine that supports full-text
search and an extensive query syntax.

>Try it out [here](https://xkcd.cheo.dev)!

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
    -n, --num       Download single comic by number
    -f, --file	    Download all comics to file
```

Start your own instance of `sxkcd` with the provided `docker-compose.yml`:

```bash
# download all comic data
$ docker compose run \
    --rm \
    --no-deps \
    app download -f /data/comics.json

$ docker-compose -d
```

### Querying
`sxkcd` supports union, negation, prefix matching and filtering by custom date ranges

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

## How it Works

`sxkcd` is a webserver built with Go and [Svelte](https://svelte.dev). It
downloads all data from xkcd and [explainxkcd](https://explainxkcd.com), before
cleaning and combining it into a single JSON file. When `sxkcd` is started, it
indexes all JSON data into Redis Stack as [JSON
documents](https://redis.io/docs/interact/search-and-query/indexing/). Redis
Stack offers native support for indexing and querying JSON documents, providing
full-text search and an extensive [query
syntax](https://redis.io/docs/interact/search-and-query/query/) for all JSON
data and its sub-elements.

Every day, `sxkcd` automatically updates the database with the newest comic, if it
doesn't already exist.

## Install

Run with docker-compose:

```bash
# download all comic data
$ docker compose run \
    --rm \
    --no-deps \
    app download -f /data/comics.json

$ docker-compose -d
```

Run locally:

```bash
# download all comic data
$ sxkcd download -f data/comics.json

# start a redis instance with data persistence
$ redis-server --appendonly yes

$ sxkcd server -p 6380 -r localhost:6379 -f data/comics.json
```

If Redis is started with persistence, `sxkcd` can be restarted without any data
files. If we wish to reindex all data in the database with a new file, we can
run `sxkcd` with the `--reindex` flag:

```bash
$ sxkcd server -p 6380 -r localhost:6379 -f data/new.json --reindex
```

This will replace all existing data with that in the new file.

## Development

`sxkcd` is built with

- Go 1.21
- [Redis Stack](https://redis.io/)
  (w/[RediSearch](https://redis.io/docs/stack/search/) and
  [RedisJSON](https://redis.io/docs/stack/json/))
- [Svelte](https://svelte.dev/) 3.46.0 and [Sveltekit](https://kit.svelte.dev/)
  1.0.0-next.405
- [picocss](https://picocss.com/) v1.5.3
- Node.js v18

### Build from Source

Clone the repository:

```bash
$ git clone https://github.com/kencx/sxkcd.git
```

The frontend's static files are embedded in the Go binary. They must generated
prior to building the Go binary:

```bash
$ cd ui && npm ci --quiet

# Build static files in ui/build
$ cd ui && npm run build

# Build binary
$ go mod download
$ make build
```

### Build Docker Image

You can also choose to build a local docker image:

```bash
$ make build
```

### Local Development

```bash
# Start frontend
$ cd ui
$ npm install
$ npm run dev

# Start redis docker container
$ docker-compose up -d redis

# Start backend
$ go mod download
$ make build
$ ./sxkcd server --port 6380 --redis localhost:6379 --file data/comics.json
```

Visit `localhost:5173`.

## Acknowledgements
>sxkcd is entirely inspired by [classes.wtf](https://github.com/ekzhang/classes.wtf)

[MIT](LICENSE)

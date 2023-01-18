#!/bin/sh

if [ -f "/data/comics.json" ]; then
	exec ./sxkcd server -p 6380 -r redis:6379 -f /data/comics.json
else
	exec ./sxkcd server -p 6380 -r redis:6379 -f https://sxkcd.s3.ap-southeast-1.amazonaws.com/comics.json
fi

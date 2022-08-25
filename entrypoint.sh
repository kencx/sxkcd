#!/bin/sh

if [ -f "comics.json" ]; then
	./rkcd server -p 6380 -r redis:6379 -f comics.json
else
	./rkcd server -p 6380 -r redis:6379 -f https://rkcd.s3.ap-southeast-1.amazonaws.com/comics.json
fi

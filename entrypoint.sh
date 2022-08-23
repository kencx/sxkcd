#!/bin/sh

if [ -f "comics.json" ]; then
	./rkcd server -p 6380 -f comics.json
else
	./rkcd server -p 6380 -f https://rkcd.s3.ap-southeast-1.amazonaws.com/comics.json
fi

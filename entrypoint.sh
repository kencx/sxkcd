#!/bin/sh

if [ -f "comics.json" ]; then
	./rkcd server -p 6380 -f comics.json
fi

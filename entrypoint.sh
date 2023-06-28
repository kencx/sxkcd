#!/bin/sh

exec ./sxkcd server -p 6380 -r redis:6379

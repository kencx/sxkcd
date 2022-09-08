#!/bin/bash
set -e

SXKCD_PATH="$HOME/sxkcd"

docker run --rm --entrypoint ./sxkcd \
	--volume "$SXKCD_PATH/app_data:/data" \
	sxkcd download -f /data/comics.json

sleep 3

if [ $(find "$SXKCD_PATH/app_data" -mtime -1 -type f -name "comics.json" 2>/dev/null) ];
then
	echo "Restarting sxkcd..."
	docker compose -f "$SXKCD_PATH/docker-compose.yml" restart
fi


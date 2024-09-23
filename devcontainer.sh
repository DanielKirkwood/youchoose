#!/bin/bash

if [ "$1" = "--up" ]; then
    docker compose -f ./.devcontainer/docker-compose.yml up --no-start
    docker compose -f ./.devcontainer/docker-compose.yml start # ensure we are started, handle also allowed to be consumed by vscode
    docker compose -f ./.devcontainer/docker-compose.yml exec devcontainer bash
fi

if [ "$1" = "--halt" ]; then
    docker compose stop
fi

if [ "$1" = "--rebuild" ]; then
    docker compose -f ./.devcontainer/docker-compose.yml up -d --force-recreate --no-deps --build devcontainer
fi

if [ "$1" = "--destroy" ]; then
    docker compose -f ./.devcontainer/docker-compose.yml down --rmi local -v --remove-orphans
fi

[ -n "$1" -a \( "$1" = "--up" -o "$1" = "--halt" -o "$1" = "--rebuild" -o "$1" = "--destroy" \) ] \
    || { echo "usage: $0 --up | --halt | --rebuild | --destroy" >&2; exit 1; }

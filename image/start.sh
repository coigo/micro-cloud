#!/bin/sh
set -e

dockerd-entrypoint.sh &

# Espera o Docker ficar disponível
until docker info >/dev/null 2>&1; do
    sleep 1
done

exec /build
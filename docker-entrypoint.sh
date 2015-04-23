#!/bin/bash
set -e

if [ "$1" = 'aprs-dashboard' ]; then
    chown -R aprs .
    exec gosu aprs "$@"
fi

exec "$@"

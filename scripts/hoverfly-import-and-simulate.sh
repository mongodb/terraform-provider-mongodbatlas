#!/bin/bash

PROXY_PORT=${1}
PROXY_ADMIN_PORT=${2}
FILE_NAME=${3}

hoverctl start --new-target "$PROXY_PORT" --proxy-port "$PROXY_PORT" --admin-port "$PROXY_ADMIN_PORT"
hoverctl import -t "$PROXY_PORT" "$FILE_NAME"
hoverctl mode -t "$PROXY_PORT" simulate

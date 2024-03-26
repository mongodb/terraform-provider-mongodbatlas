#!/bin/bash

PROXY_PORT=${1}
PROXY_ADMIN_PORT=${2}

hoverctl start --new-target "$PROXY_PORT" --proxy-port "$PROXY_PORT" --admin-port "$PROXY_ADMIN_PORT"
hoverctl mode -t "$PROXY_PORT" capture --stateful

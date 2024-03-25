#!/bin/bash

PROXY_PORT=${1}
FILE_NAME=${2}

hoverctl -t $PROXY_PORT export $FILE_NAME
hoverctl stop --target $PROXY_PORT
hoverctl targets delete $PROXY_PORT --force

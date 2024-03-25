#!/bin/bash

PROXY_PORT=${1}

hoverctl stop --target $PROXY_PORT
hoverctl targets delete $PROXY_PORT --force

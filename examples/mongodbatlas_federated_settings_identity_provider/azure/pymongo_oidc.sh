#!/bin/bash
set -euo pipefail

sudo apt-get install -y python3-pip
pip install pymongo

export DATABASE="${DATABASE}"
export COLLECTION="${COLLECTION}"
export RECORD='${RECORD}' # single quotes for json payload
export MONGODB_URI="${MONGODB_URI}"

sudo chown ${OS_USER}:${OS_USER} /home/${OS_USER}/pymongo_oidc.py
python3 /home/${OS_USER}/pymongo_oidc.py

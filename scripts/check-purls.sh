#!/usr/bin/env bash
set -euo pipefail

if ! git diff --quiet --exit-code compliance/purls.txt; then
    echo "compliance/purls.txt is out of date. Please run 'make gen-purls' and commit the result."
    git --no-pager diff compliance/purls.txt
    exit 1
fi
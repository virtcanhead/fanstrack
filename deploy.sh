#!/bin/bash

set -e
set -u

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

docker build -t canhead/fanstrack .
docker push canhead/fanstrack
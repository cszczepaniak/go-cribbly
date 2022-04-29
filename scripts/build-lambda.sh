#!/bin/bash

set -ex

source ./scripts/variables.sh

GOOS=linux CGO_ENABLED=0 go build -o cribbly-backend cmd/lambda/main.go
zip $BUILD_HASH.zip cribbly-backend
rm cribbly-backend
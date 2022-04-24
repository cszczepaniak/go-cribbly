#!/bin/bash

set -ex

source ./scripts/variables.sh

cd infrastructure/
cdk deploy $INFRASTRUCTURE_STACK --require-approval=never
cd -

GOOS=linux CGO_ENABLED=0 go build -o cribbly-backend cmd/lambda/main.go
zip $BUILD_HASH.zip cribbly-backend
rm cribbly-backend

aws s3 cp $BUILD_HASH.zip s3://$CRIBBLY_APP_BUCKET/
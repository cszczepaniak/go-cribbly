#!/bin/bash

set -ex

source ./scripts/variables.sh

cd infrastructure/
cdk deploy $INFRASTRUCTURE_STACK --require-approval=never
cd -

aws s3 cp $BUILD_HASH.zip s3://$CRIBBLY_APP_BUCKET/
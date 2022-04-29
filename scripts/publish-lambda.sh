#!/bin/bash

set -ex

source ./scripts/variables.sh

aws s3 cp $BUILD_HASH.zip s3://$CRIBBLY_APP_BUCKET/
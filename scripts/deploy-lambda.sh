#!/bin/bash

set -ex

source ./scripts/variables.sh

cd infrastructure/
cdk deploy $APP_STACK --require-approval=never
cd -
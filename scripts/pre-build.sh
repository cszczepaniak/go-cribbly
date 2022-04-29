#!/bin/bash

set -ex

source ./scripts/variables.sh

cd infrastructure/
cdk deploy $INFRASTRUCTURE_STACK --require-approval=never
cd -
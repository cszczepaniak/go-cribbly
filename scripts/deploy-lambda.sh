#!/bin/bash

set -ex

source ./scripts/variables.sh

cdk deploy $APP_STACK --require-approval=never
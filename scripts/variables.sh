#!/bin/bash

export BUILD_HASH=${GITHUB_SHA:-"unversioned"}
export CRIBBLY_APP_BUCKET="cribbly-app-bucket"
export INFRASTRUCTURE_STACK="CribblyInfrastructureStack"
export APP_STACK="CribblyApplicationStack"
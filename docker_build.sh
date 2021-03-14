#!/bin/bash

############################################################################################
#
# This script builds golang app using Docker container
#
############################################################################################

IMAGE="golang"
TAG="1.16"
GO_COMMAND="go build -v"

sudo docker run --rm -v ${PWD}:/usr/src/myapp -w /usr/src/myapp ${IMAGE}:${TAG} ${GO_COMMAND}

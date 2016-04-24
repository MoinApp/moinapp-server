#!/bin/bash

BRANCH=$(git symbolic-ref --short -q HEAD)
if [ $BRANCH = "" ];
  then
  $BRANCH = git rev-parse HEAD
fi

go build -ldflags "-X main.APP_VERSION=$BRANCH"

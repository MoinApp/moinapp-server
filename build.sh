#!/bin/bash

if [[ $TRAVIS = true ]] && [[ $CI = true ]]; then
  # if running on Travis CI

  if [[ $TRAVIS_TAG ]]; then
    # if building a tag
    BRANCH=$TRAVIS_TAG
  else
    # else we are building a branch
    BRANCH=$TRAVIS_BRANCH
  fi
else
  # anywhere else, ask git!
  BRANCH=$(git symbolic-ref --short -q HEAD)
  # or if our head is detached, ask for the commit hash
  if [[ $BRANCH = "" ]];
    then
    BRANCH = git rev-parse HEAD
  fi
fi

go build -ldflags "-X main.APP_VERSION=$BRANCH"

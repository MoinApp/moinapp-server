#!/bin/bash

if [[ $TRAVIS = true ]] && [[ $CI = true ]]; then
  # if running on Travis CI

  if [[ $TRAVIS_TAG ]]; then
    # if building a tag
    NAME=$TRAVIS_TAG
  else
    # else we are building a commit
    NAME=$TRAVIS_COMMIT
  fi
else
  # anywhere else, ask git!
  NAME=$(git rev-parse HEAD)
fi

go build -ldflags "-X github.com/MoinApp/moinapp-server/info.AppVersion=$NAME"

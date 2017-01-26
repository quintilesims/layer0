#!/bin/bash
# Travis CI Release / Smoketest Script
#
LAYER0_PATH=$GOPATH/src/github.com/quintilesims/layer0
LAYER0_PREFIX=smoketest
set -e
set -x

# Ignore if master branch is not target
if [ ! "$TRAVIS_BRANCH" == "master" ]; then
  echo "[INFO] Skipping Smoketest";
  exit 0
fi

# Tag repository
git tag "$LAYER0_PREFIX"
env

# Install and apply smoketest
make install-smoketest
make apply-smoketest

# export l0 environment
cd "$LAYER0_PATH"/setup
eval "$(./l0-setup endpoint -i "$LAYER0_PREFIX")"
cd -

# run smoketest
make smoketest
echo "[INFO] Smoketest completed successfully!"

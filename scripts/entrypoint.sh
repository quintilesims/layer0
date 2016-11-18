#!/bin/sh

config="{ \"auths\": { \"d.ims.io\": { \"auth\": \"$DOCKER_TOKEN\", \"email\": \"\" }}}"
mkdir -p /root/.docker
echo $config > /root/.docker/config.json

exec "$@"

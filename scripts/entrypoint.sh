#!/bin/sh

config="{ \"auths\": { \"hub.docker.com\": { \"auth\": \"$DOCKER_TOKEN\", \"email\": \"\" }}}"
mkdir -p /root/.docker
echo $config > /root/.docker/config.json

exec "$@"

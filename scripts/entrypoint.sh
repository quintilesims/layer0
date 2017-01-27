#!/bin/sh

config="{ \"auths\": { \"hub.docker.com\": { \"auth\": \"$DOCKER_TOKEN\", \"email\": \"\" }}}"
mkdir -p /root/.docker
echo $config > /root/.docker/config.json

# start and configure MySql on boot
service mysql start
mysql -u root -e "CREATE USER 'layer0'@'127.0.0.1' IDENTIFIED BY 'nohaxplz';"
mysql -u root -e "GRANT ALL ON *.* TO 'layer0'@'127.0.0.1';"

exec "$@"

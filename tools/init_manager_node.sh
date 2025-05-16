#!/bin/bash

exec 2>&1

# Base packages installation

apt-get update
apt-get install net-tools

# PostgreSQL SSL cert installation

mkdir -p ./.data/yc && \
wget "https://storage.yandexcloud.net/cloud-certs/CA.pem" \
    --output-document ./.data/yc/root.crt && \
chmod 0655 ./.data/yc/root.crt

# Certbot directories creation

mkdir -p ./.data/certbot/conf
mkdir -p ./.data/certbot/webroot

# Docker installation & configuration

apt-get install docker.io

DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
mkdir -p $DOCKER_CONFIG/cli-plugins
curl -o $DOCKER_CONFIG/cli-plugins/docker-compose \
    -SL https://github.com/docker/compose/releases/download/v2.32.4/docker-compose-linux-x86_64

chmod +x $DOCKER_CONFIG/cli-plugins/docker-compose
chmod 0666 /var/run/docker.sock

echo $YANDEX_OAUTH_TOKEN | docker login \
    --username oauth \
    --password-stdin \
    cr.yandex

docker swarm init

#!/bin/bash

exec 2>&1

# Base packages installation

apt-get update
apt-get install net-tools

# Yandex Cloud SSL CA certificate installation

mkdir -p ./.data/yc && \
wget "https://storage.yandexcloud.net/cloud-certs/CA.pem" \
    --output-document ./.data/yc/root.crt && \
chmod 0655 ./.data/yc/root.crt

# Docker installation & configuration

apt-get install docker.io

chmod 0666 /var/run/docker.sock

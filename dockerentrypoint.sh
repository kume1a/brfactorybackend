#!/bin/sh

echo "172.17.0.1 host.docker.internal" >> /etc/hosts

exec ./app serve --http=0.0.0.0:8090
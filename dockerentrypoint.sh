#!/bin/sh

# Update /etc/hosts
echo "172.17.0.1 host.docker.internal" >> /etc/hosts

# Start your app
exec ./app serve --http=0.0.0.0:8090
#!/bin/zsh
# Manually copying dependencies is ugly, but it's the only way I've found to copy dependencies locally
mkdir -p ./vendor/github.com/svera/acquire && cp -r ../acquire ./vendor/github.com/svera
docker-compose build 
docker-compose up
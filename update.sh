#!/bin/bash

docker pull $1/prod:latest
if docker stop app; then docker rm app; fi
docker run -d -p 8080:8080 --name ma-app $1/prod
if docker rmi $(docker images --filter "dangling=true" -q --no-trunc); then :; fi

#!/bin/sh
configPath=$(pwd)/config/
cd ./sqlgpt
cligptDir=$(pwd)
docker container rm sqlgpt
docker image build -t sqlgpt:latest .
docker run --rm -d --name sqlgpt --network sqlgpt -p 8080:8080 -v $configPath:/config sqlgpt

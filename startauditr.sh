#!/bin/bash
configPath=$(pwd)/config/

cd ./auditr
docker build -t auditr .
docker run --rm -d --name auditr --network sqlgpt  --volume $configPath:/config -v /etc/ssl/certs:/etc/ssl/certs:ro auditr

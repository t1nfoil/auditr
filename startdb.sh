#!/bin/bash

cd database
currentDir=$(pwd)
docker build -t my-postgres .

docker run -d --rm --name postgres-dev --network sqlgpt \
  --volume $currentDir/pgdata:/var/lib/pgsql/data \
  --volume $currentDir/mnt_data:/mnt/data \
  --volume $currentDir/pg_hba.conf:/etc/postgresql/pg_hba.conf \
  --volume $currentDir/postgresql.conf:/etc/postgresql/postgresql.conf \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_USER=postgres \
  -e PGDATA=/var/lib/pgsql/data/pgdata14 \
  -e POSTGRES_INITDB_ARGS="--data-checksums --encoding=UTF8" \
  -e POSTGRES_DB=db \
  -p 5432:5432 \
  my-postgres \
  postgres \
    -c 'config_file=/etc/postgresql/postgresql.conf' \
    -c 'hba_file=/etc/postgresql/pg_hba.conf'

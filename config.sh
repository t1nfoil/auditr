#!/bin/bash
export PGPASSWORD='password'

docker network create sqlgpt

# start the database
./startdb.sh

#import the database tables
psql -U postgres -h localhost -f ./database/create_database.sql
psql -U postgres -h localhost -f ./database/create_tables_da.sql
psql -U postgres -h localhost -f ./database/create_tables_si.sql

#start sqlgpt
./startsqlgpt.sh

#start auditr
./startauditr.sh


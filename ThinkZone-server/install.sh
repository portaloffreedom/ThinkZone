#!/bin/bash

initdb -U thinkzone -D thinkzoneDB
postgres -D . >postgres.log 2>&1 &
createdb thinkzoneDB
psql -d thinkzoneDB -f ./src/thinkzone/database/create_database.sql
psql -d thinkzoneDB -c "INSERT INTO conversation VALUES (0)"
psql -d thinkzoneDB -c "INSERT INTO post VALUES (0,0,null,-1,-1,-1)"

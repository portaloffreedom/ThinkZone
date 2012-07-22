#!/bin/bash

initdb -U thinkzone -D thinkzoneDB
#postgres -D thinkzoneDB >postgres.log 2>&1 &
pg_ctl -D thinkzoneDB -l postres.log start
createuser thinkzone -s 
createdb --owner=thinkzone thinkzoneDB
psql -d thinkzoneDB -f src/thinkzone/database/create_database.sql
psql -d thinkzoneDB -c "INSERT INTO conversation VALUES (0)"
psql -d thinkzoneDB -c "INSERT INTO post VALUES (0,0,'senza titolo',-1,-1,-1)"

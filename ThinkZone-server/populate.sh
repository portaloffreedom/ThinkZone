#!/bin/bash
psql -d thinkzoneDB -f src/thinkzone/database/create_database.sql
psql -d thinkzoneDB -c "INSERT INTO conversation VALUES (0)"
psql -d thinkzoneDB -c "INSERT INTO post VALUES (0,0,'senza titolo',-1,-1,-1)"

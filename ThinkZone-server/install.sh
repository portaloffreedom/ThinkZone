#!/bin/bash

# Parte iniziale tolta petché postgres mi fa "girare i gangheri"
#initdb -U thinkzone -D thinkzoneDB
####postgres -D thinkzoneDB >postgres.log 2>&1 &
#pg_ctl -D thinkzoneDB -l postres.log start
echo "serve il diritto di editare il database sennò non si riesce a fare l'init del database"

# Creazione del database
createuser thinkzone -s 
createdb --owner=thinkzone thinkzoneDB

# Riempimento del database con le tabelle necessarie
sh populate.sh

echo
echo "init del database finito"
echo
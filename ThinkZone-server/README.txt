Come installare questo programma:

**********************
**** Compilazione ****
**********************
# Dipendenze:

  go - http://golang.org/
    seguire il link per le istruzioni di come installare:
    http://golang.org/doc/install

  go postresql driver - https://github.com/jbarham/gopgsqldriver
    eseguire il seguente comando per installare il driver (probabilmente dovrà essere 
    eseguito con i privilegi di amministratore)
      $ go get github.com/jbarham/gopgsqldriver
    Secondo quello che afferma l'autore di questo driver l'installazione è configurata 
    per funzionare su linux - potrebbe dare problemi su altri sistemi

# Comandi da eseguire:

sh build.sh

***********************
**** Installazione ****
***********************
Questo programma non ha necessità di essere installato. Funziona benissimo nella
cartella dove è stato scaricato

# Servono i seguenti programmi installati per potere avviare il server:
  
  postresql - http://www.postgresql.org/
    seguire il link per le istruzioni di come installare:
    http://www.postgresql.org/download/

***********************
******** Avvio ********
***********************
Requisisti per il primo avvio del server

  server posgresql deve essere avviato prima di fare partire il server
    su ubuntu dovrebbe essere il seguente comando:
    $ sudo service postgresql start
  
  Se è il primo avvio va eseguito lo script per popolare il server:
    $ sh populate.sh
    
  Il programma è pronto per essere avviato:
    $ ./ThinkZoneServer
  
  
  
  
  
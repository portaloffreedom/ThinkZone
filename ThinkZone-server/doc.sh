#!/bin/bash

helptext="Comandi disponibili:\n"\
"\t help\t mostra questo aiuto\n"\
"\t kill\t uccide il server http che tiene occupata la porta 8080"

# nessun comando: azione di default
if [ $# -eq 0 ]
  then 	echo "DOCUMENTAZIONE"
	echo "Assicurarsi di avere la porta 8080 libera per fare partire il server"\
	     " della documentazione"
	
	# AVVIO DEL SERVER DELLA DOCUMENTAZIONE
	godocpid=`pidof godoc`
	if [ $? -ne 0 ]
	  then	GOPATH=/home/matteo/progetti/thinkzoneserver godoc -http=:8080 &
	fi
	
	# AVVIO DEL BROWSER SULLA DOCUMENTAZIONE DI THINKZONE
	xdg-open http://localhost:8080/pkg/thinkzone/
	if [ $? -ne 0 ] 
	  then	echo "xdg non installato o non funzionante correttamente. "\
		     "Aprire manualmente il browser sulla pagina:"
		echo "http://localhost:8080/pkg/thinkzone/"
	fi
	
	# STAMPA MEMO
	echo "ricordarsi di fare il comando \"sh doc.sh kill\" per "\
	     "terminare il server http in locale e liberare la porta 8080 "\
	     "sul proprio pc"
	     
  else 	# comando help
	if [ $1 == "help" ]
	  then	echo -e $helptext
	fi

	# comando kill
	if [ $1 == "kill" ]
	  then	killall godoc
		if [ $? -eq 0 ]
		  then	echo "ucciso con successo il server della documentazione"
		  else	echo "server della documentazione non ucciso"
		fi  
	fi
fi


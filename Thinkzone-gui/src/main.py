'''
Classe main per il programma Thinkzone.
Permette di registrare un utente via riga di comando, o eventualmente di far partire l'interfaccia grafica.

Utilizzo: Thinkzone [-h] [-v] (-g | -r username password) [-a hostname]
                 [-p porta]

Argomenti opzionali:
  -h, --help            Visualizza questo messaggio.
  -v, --version         stampa la versione del programma ed esce.
  -g, --gui             Fa partire l'interfaccia grafica.
  -r username password  Registra un utente.

Opzioni registrazione:
  -a hostname, --address hostname
                        Hostname del server per la registraizone.
  -p porta, --port porta
                        La porta del server a cui connettersi.


@author: stengun
'''
import argparse,sys
from rete import Comunicazione
from gui import MainWindow
from PyQt4 import QtGui

__VERSION__ = "0.0.9"

if __name__ == '__main__':
    argparser = argparse.ArgumentParser(prog="Thinkzone",version = __VERSION__)
    grupporeg = argparser.add_argument_group("Opzioni registrazione")
    gruppoescl = argparser.add_mutually_exclusive_group(required = True)
    gruppoescl.add_argument("-g","--gui",action="store_true",help="Fa partire l'interfaccia grafica.")
    gruppoescl.add_argument("-r",nargs=2,metavar=("username","password"),help="Registra un utente.")
    grupporeg.add_argument("-a","--address",metavar="hostname",help="Hostname del server per la registraizone.")
    grupporeg.add_argument("-p","--port",metavar="porta",type=int,default=4242,help="La porta del server a cui connettersi.")
    arogmi = argparser.parse_args(sys.argv[1:])
    if(arogmi.gui):
        print('Starting gui...')
        app = QtGui.QApplication(sys.argv)
        finestra = MainWindow.mainwindow(__VERSION__)
        finestra.show()
        app.exec()
        print('exit')
    else:
        username = arogmi.r[0]
        password = arogmi.r[1]
        porta = arogmi.port
        hostname = arogmi.address
        connettore = Comunicazione.comunicatore()
        connettore.registrati(hostname, porta, username, password)
        if(connettore._registered):
            print('Registrazione completata correttamente')
        else:
            print('Registrazione non effettuata. Consulta il file di log per ulteriori informazioni.')
'''
Created on 11/lug/2012
Classe comunicatore che riesce a ricevere e inviare dati.
Ha dentro definite le meccaniche di comunicazione.
@author: stengun
'''
import queue
import socket
import sys
from PyQt4 import QtCore

class comunicatore(QtCore.QThread):
    
    _socket = None
    _posizione = None
    _messaggi = None
    _stop = None
    cursore = None
    _utenteAttivo = None
    _userID = None
    _receive_thread = None
    
    def __init__(self):
        QtCore.QThread.__init__(self)
        self._socket = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        self._messaggi = queue.Queue(255)
        self._stop = False
        rimozione = QtCore.pyqtSignal(int,int,name='rimozione')
        aggiunta = QtCore.pyqtSignal(int,str,name='aggiunta')
        
    def run(self):
        while(not(self._stop)):
            messaggio = self._messaggi.get(True, None)
            if(self._parseinput(messaggio)):
                print('emetto')
                self.emit(QtCore.SIGNAL('aggiunta(int,str)'),self.cursore,messaggio)
            
    
    def _parseinput(self,messaggio):
        if(messaggio == '\\'):
            messaggio = self._messaggi.get(False,None)
        else:
            return True
        if(messaggio == '\\'):
            return True
        if(messaggio == 'P'):
            self.cursore = self._messaggi.get(False,None)
            controllo = self._messaggi.get(False,None)
            if(controllo != '\\'):
                print('Errore stream TCP!',file=sys.stderr)
                self.cursore = None
            return False
        if(messaggio == 'D'):
            quantita = self._messaggi.get(False,None)
            self.cursore -= quantita
            controllo = self._messaggi.get(False,None)
            if(controllo != '\\'):
                print('Errore stream TCP!',file=sys.stderr)
                self.cursore = None
            else:
                if(self._utenteAttivo != self._userID):
                    self.emit(self,QtCore.SIGNAL('rimozione(int,int)'),self._cursore,quantita)
            return False
        if(messaggio== 'U'):
            self._utenteAttivo = self._messaggi.get(False,None)
            controllo = self._messaggi.get(False,None)
            if(controllo != '\\'):
                print('Errore stream TCP!',file=sys.stderr)
                self.cursore = None
            return False
        
    def connetti(self,hostname,porta,nickname):
        self._socket.connect((hostname,porta))
        self._receive_thread = Receiver(self._messaggi,self._socket)
        self._receive_thread.start()
        self._spedisci(nickname+'\\')
        messaggio = self._messaggi.get(True, None)
        self._userID = messaggio
        
    def disconnetti(self):
        self._socket.close()
        
    def _spedisci(self,data):
        self._socket.sendall(data.encode())
        
    def spedisci_aggiunta(self,posizione,dati):
        #dati = dati.encode()
        if(posizione != self._posizione):
            self._spedisci('\P'+str(posizione)+'\\')
        self._spedisci(dati)
        self._posizione = posizione+1
        
    def spedisci_rimozione(self,posizione,rimossi):
        posizione += rimossi
        if(posizione != self._posizione):
            self._spedisci('\P'+str(posizione)+'\\')
        self._spedisci('\D'+str(rimossi)+'\\')
        self._posizione = posizione-1
        
class Receiver(QtCore.QThread):
    
    _coda = queue.Queue()
    _stop = None
    _socket = None
    def __init__(self,coda,socket):
        QtCore.QThread.__init__(self)
        self._coda = coda
        self._stop = False
        self._socket = socket
    
    def run(self):
        while(not(self._stop)):
            print('ricevimento')
            self._coda.put(self._socket.recv(1))
    
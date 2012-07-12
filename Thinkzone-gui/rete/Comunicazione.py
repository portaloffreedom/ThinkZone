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
    blink_cursor = None
    _utenteAttivo = None
    _userID = None
    _receive_thread = None
    
    def __init__(self):
        QtCore.QThread.__init__(self)
        self._socket = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        self._messaggi = queue.Queue(255)
        self._stop = False
        _rimozione = QtCore.pyqtSignal(int,int,name='rimozione')
        _aggiunta = QtCore.pyqtSignal(int,str,name='aggiunta')
        self.blink_cursor = 0
        
    def run(self):
        while(not(self._stop)):
            messaggio = self._messaggi.get(True, None)
            messaggio = (messaggio).decode("utf-8")
            #print('messaggio',messaggio)
            if(self._parseinput(messaggio) and self._utenteAttivo != self._userID):
                print('emetto',messaggio)
                self.emit(QtCore.SIGNAL('aggiunta(int,QString)'),self.blink_cursor,messaggio)
                self.blink_cursor += 1
           
    def _controller(self,controllo,messaggio):
        if(controllo != '\\'):
            print('Errore stream TCP!',file=sys.stderr)
            self.blink_cursor = 0 
            messaggio = None
        return messaggio
    
    def _recvInt(self):
        intero = 0
        next = '\00'
        while(True):
            next = self._messaggi.get(True,None).decode("utf-8")
            if(next == '\\' or int(next) < 0 or int(next) > 9):
                return intero,next
            intero = intero*10+int(next)
    
    def _parseinput(self,messaggio):
        controllo = '\00'
        if(messaggio == '\\'):
            print('caso 1')
            messaggio = self._messaggi.get(True,None).decode("utf-8")
            if(messaggio == '\\'):
                return True
            if(messaggio == 'P'):
                print('caso P')
                [self.blink_cursor, controllo] = self._recvInt()
                messaggio = self._controller(controllo, messaggio)
                return False
            if(messaggio == 'D'):
                print('caso D')
                [quantita, controllo] = self._recvInt()
                self.blink_cursor -= quantita
                #controllo = self._messaggi.get(True,None)
                messaggio = self._controller(controllo, messaggio)
                if(messaggio == None):
                    return False
                else:
                    if(self._utenteAttivo != self._userID):
                        self.emit(QtCore.SIGNAL('rimozione(int,int)'),self.blink_cursor,quantita)   
                return False
            if(messaggio== 'U'):
                print('caso U')
                [self._utenteAttivo,controllo] = self._recvInt()
                print('utente attivo',self._utenteAttivo)
                if(self._utenteAttivo != self._userID):
                    self.cursore_locale = self.blink_cursor
                messaggio = self._controller(controllo, messaggio)
                return False
        return True
    
    def connetti(self,hostname,porta,nickname):
        self._socket.connect((hostname,porta))
        self._receive_thread = Receiver(self._messaggi,self._socket)
        self._receive_thread.start()
        self._spedisci(nickname+'\\')
        controllo = '\00'
        [self._userID,controllo] = self._recvInt()
        print('il tuo ID',self._userID)
        print(controllo)
        
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
    
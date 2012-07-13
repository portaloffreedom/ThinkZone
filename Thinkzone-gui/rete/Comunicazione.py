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
    _response = None
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
            try:
                messaggio = self._messaggi.get(True, None)
            except:
                print('errore connettore',sys.exc_info())
                self._stop = True
                self._receive_thread.setTerminationEnabled()
                self._receive_thread._stop = True
                continue
            print('messaggio',messaggio)
            if(self._parseinput(messaggio) and self._utenteAttivo != self._userID):
                print('emetto',messaggio)
                self.emit(QtCore.SIGNAL('aggiunta(int,QString)'),self.blink_cursor,messaggio)
                self.blink_cursor += 1
        try:
            self._socket.close()
        except:
            print('problema',sys.exc_info())
        finally:
            self._stop = False
            self._socket = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
           
    def _controller(self,controllo,messaggio):
        if(controllo != '\\'):
            print('Errore stream TCP!',file=sys.stderr)
            self.blink_cursor = 0 
            messaggio = None
        return messaggio
    
    def _recvInt(self):
        intero = 0
        prossimo = '\00'
        while(True):
            prossimo = self._messaggi.get(True,None)
            if(prossimo == '\\' or int(prossimo) < 0 or int(prossimo) > 9):
                return intero,prossimo
            intero = intero*10+int(prossimo)
    
    def _parseinput(self,messaggio):
        controllo = '\00'
        if(messaggio == '\\'):
            print('caso 1')
            messaggio = self._messaggi.get(True,None)
            print('comando:',messaggio)
            if(messaggio == '\\'):
                return True
            if(messaggio == 'C'):
                print('caso C')
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
            if(messaggio =='R'):
                print('Caso Response')
                [self._response,controllo] = self._recvInt()
                print('risposta',self._response)
                messaggio = self._controller(controllo, messaggio)
                print(controllo)
                return False
        return True
    
    def registrati(self,hostname,porta,nickname,password):
        self._socket.connect((hostname,porta))
        self._receive_thread = Receiver(self._messaggi,self._socket)
        self._receive_thread.start()
        self._spedisci('\L1\\')
        self._spedisci(nickname+'\\'+password+'\\')
        messaggio = self._messaggi.get(True, None)
        if(self._parseinput(messaggio) or self._parseResponse(self._response)):
            print('Errore di registrazione! errore ',self._response)
            sys.exit()
        print('Registrazione completata!')
        self._socket.close()
        self._socket = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        self._receive_thread.setTerminationEnabled(True)
        self._receive_thread._stop = True
        self._messaggi.put('close')
        print(self._messaggi.get())
    
    def connetti(self,hostname,porta,nickname,password):
        self._socket.connect((hostname,porta))
        self._receive_thread = Receiver(self._messaggi,self._socket)
        self._receive_thread.start()
        
        self._spedisci('\L0\\')
        self._spedisci(nickname+'\\'+password+'\\')
        messaggio = self._messaggi.get(True, None)
        
        if(self._parseinput(messaggio)):
            print('Errore di login! Errore',self._response)
            return
        
        controllo = '\00'
        if(self._parseResponse(self._response)):
            errore = 'il server non ha accettato la connessione'
            try:
                self._receive_thread._stop = True
                self._socket.close()
            except:
                errore = 'il server ha chiuso la connessione!'+sys.exc_info()
            finally:
                print(errore,file=sys.stderr)
                self._stop = True
                self._messaggi.put('close')
                print(self._messaggi.get())
                return
        try:
            [self._userID,controllo] = self._recvInt()
        except:
            print('connessione chiusa!')
            sys.exit()
        print('il tuo ID',self._userID)
        print(controllo)
        
    def _parseResponse(self,response):
        if(response == 0):
            return False
        return True
    
    def disconnetti(self):
        self._socket.close()
        
    def _spedisci(self,data):
        try:
            self._socket.sendall(data.encode())
        except:
            print('error',sys.exc_info())
        
    def spedisci_aggiunta(self,posizione,dati):
        #dati = dati.encode()
        if(posizione != self._posizione):
            self._spedisci('\C'+str(posizione)+'\\')
        self._spedisci(dati)
        self._posizione = posizione+1
        
    def spedisci_rimozione(self,posizione,rimossi):
        posizione += rimossi
        if(posizione != self._posizione):
            self._spedisci('\C'+str(posizione)+'\\')
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
            try:
                buffer = self._socket.recv(1)
                try:
                    buffer = buffer.decode("utf-8")
                except:
                    buffer += self._socket.recv(1)
                    buffer = buffer.decode("utf-8")
                print(buffer)
                self._coda.put(buffer)
            except:
                print('eccezione',sys.exc_info())
                self._stop = True
        self._stop = False
    
'''
Classe comunicatore che riesce a ricevere e inviare dati.
Ha dentro definite le meccaniche di comunicazione.
@author: stengun
'''
import queue
import socket
import sys
from PyQt4 import QtCore

class comunicatore(QtCore.QThread):
    '''
    Generico comunicatore che imposta anche i thread per la ricezione.
    '''
    _socket = None
    _posizione = None
    _messaggi = None
    _stop = None
    blink_cursor = None
    _utenteAttivo = None
    _userID = None
    _receive_thread = None
    _response = None
    _activePost = None
    _barrier = None
    _cursors = {}
    def __init__(self):
        QtCore.QThread.__init__(self)
        self._socket = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        self._messaggi = queue.Queue(255)
        self._stop = False
        _rimozione = QtCore.pyqtSignal(int,int,name='rimozione')
        _aggiunta = QtCore.pyqtSignal(int,str,name='aggiunta')
        _nuovopost = QtCore.pyqtSignal(int,name='nuovoPost')
        _selpost = QtCore.pyqtSignal(int,name='selectPost')
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
            #print('messaggio',messaggio)
            if(self._parseinput(messaggio) and self._utenteAttivo != self._userID):
                print('scrivo:',messaggio,
                      'per utente',self._utenteAttivo,
                      'con cursore',self._cursors[self._utenteAttivo][0],
                      'e post',self._cursors[self._utenteAttivo][1])
                
                self.emit(QtCore.SIGNAL('aggiunta(int,QString)'),self._cursors[self._utenteAttivo][0],messaggio)
                self._cursors[self._utenteAttivo] = (self._cursors[self._utenteAttivo][0]+1,self._cursors[self._utenteAttivo][1])
                
        try:
            self.disconnetti()
        except:
            print('problema',sys.exc_info())
        finally:
            self._stop = False
           
    def _controller(self,controllo,messaggio):
        '''
        Questo metodo decide se il carattere "controllo" è un carattere di fine comando.
        Se questo test fallisce, imposta il valore di messaggio a None.
        '''
        if(controllo != '\\'):
            print('Errore stream TCP!',file=sys.stderr)
            self.blink_cursor = 0 
            messaggio = None
        return messaggio
    
    def _recvInt(self):
        '''
        Questo metodo riceve un numero intero da uno stream di input.
        Ritorna due valori: [numero intero, carattere dopo l'intero]
        '''
        intero = 0
        prossimo = '\00'
        while(True):
            prossimo = self._messaggi.get(True,None)
            if(prossimo == '\\' or int(prossimo) < 0 or int(prossimo) > 9):
                return intero,prossimo
            intero = intero*10+int(prossimo)
    
    def _parseinput(self,messaggio):
        '''
        Effettua l'analisi del messaggio ricevuto: Ritorna TRUE se si tratta di un carattere di controllo,
        altri menti ritorna False quando si tratta di un comando.
        '''
        controllo = '\00'
        
        if(messaggio == '\\'): # selezione comando
            #print('caso 1')
            messaggio = self._messaggi.get(True,None)
            #print('comando:',messaggio)
            if(messaggio == '\\'): #non è un comando.
                return True
            
            if(messaggio == 'C'): #Modifica cursore
                print('caso C')
                [self.blink_cursor, controllo] = self._recvInt()
                print('cursore ',self.blink_cursor)
                self._cursors[self._utenteAttivo] = (self.blink_cursor,self._activePost)
                messaggio = self._controller(controllo, messaggio)
                return False
            
            if(messaggio == 'D'): # Eliminazione
                print('caso D')
                [quantita, controllo] = self._recvInt()
                self.blink_cursor -= quantita
                self._cursors[self._utenteAttivo] = (self.blink_cursor,self._activePost)
                messaggio = self._controller(controllo, messaggio)
                if(messaggio == None):
                    return False
                else:
                    if(self._utenteAttivo != self._userID):
                        self.emit(QtCore.SIGNAL('rimozione(int,int)'),self._cursors[self._utenteAttivo][0],quantita)   
                return False
            
            if(messaggio== 'U'): # Selezione utente
                #print('caso U')
                [self._utenteAttivo,controllo] = self._recvInt()
                print('utente attivo',self._utenteAttivo)
                try:
                    self._cursors[self._utenteAttivo]
                except:
                    print('Warning: utente nuovo',file=sys.stderr)
                    self._cursors[self._utenteAttivo] = (0,0)

                #if(self._utenteAttivo != self._userID):
                #    self.cursore_locale = self.blink_cursor
                messaggio = self._controller(controllo, messaggio)
                return False
            
            if(messaggio =='R'): #Risposta a una azione
                #print('Caso Response')
                [self._response,controllo] = self._recvInt()
                print('Caso Response: risposta',self._response)
                messaggio = self._controller(controllo, messaggio)
                #print(controllo)
                return False
            
            if(messaggio == 'P'):#creazione post
                print('caso POST')
                [idpost,controllo] = self._recvInt()
                #print('idPost',idpost)
                messaggio = self._controller(controllo, messaggio)
                #self._activePost = idpost
                actu = 0
                try:
                    temp = self._cursors[self._utenteAttivo]
                    if(temp[1] == idpost):
                        actu = temp[0]
                finally:
                    self._cursors[self._utenteAttivo] = (actu,idpost)
                    
                self.emit(QtCore.SIGNAL('selectPost(int)'),idpost)
                self._barrier.wait()
                return False
            
            if(messaggio == 'K'):
                print('caso Kreazione') #non sono un bimbominchia
                [parent,controllo] = self._recvInt()
                [idpost,controllo] = self._recvInt()
                self.emit(QtCore.SIGNAL('nuovoPost(int)'),idpost)
                self._barrier.wait()
                return False
            
        return True
    
    def registrati(self,hostname,porta,nickname,password):
        '''
        Metodo per effettuare una registrazione ad un server.
        '''
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
        self.disconnetti()
    
    def connetti(self,hostname,porta,nickname,password):
        '''
        Metodo per effettuare una connessione ad un server.
        '''
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
                self.disconnetti()
            except:
                errore = 'il server ha chiuso la connessione!'+sys.exc_info()
            finally:
                print(errore,file=sys.stderr)
                self._stop = True
                return
        try:
            [self._userID,controllo] = self._recvInt()
        except:
            print('connessione chiusa!',file=sys.stderr)
            sys.exit()
        print('il tuo ID',self._userID)
        print(controllo)
        self._cursors[self._userID] = (0,0)
        
    def _parseResponse(self,response):
        '''
        Metodo che analizza i codici di risposta, e ritorna True se è un codice di errore.
        '''
        if(response == 0):
            return False
        return True
    
    def disconnetti(self):
        '''
        Classe che ti disconnette da un server.
        Uccide il thread che ascolta e reimpostaa zero tutte le strutture dati.
        '''
        self._socket.shutdown(socket.SHUT_RDWR)
        self._socket.close()
        self._socket = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        self._receive_thread.setTerminationEnabled(True)
        self._receive_thread._stop = True
        self._messaggi = queue.Queue(255)
        
    def _spedisci(self,data):
        '''
        Spedisce dei dati al server
        '''
        try:
            self._socket.sendall(data.encode())
        except:
            print('errore di spedizione',sys.exc_info())
        
    def spedisci_aggiunta(self,posizione,dati,idpost):
        '''
        Spedisce al server una aggiunta di testo su uno specifico post e da una specifica posizione.
        '''
        #dati = dati.encode()
        print('posizione attuale',self._cursors[self._utenteAttivo][0])
        print('posizione da cui partire',posizione)
        print('posizione da aggiungere',len(dati))
        
        if(self._cursors[self._utenteAttivo][1] != idpost):
            self._spedisci('\P'+str(idpost)+'\\')
            #self._cursors[self._utenteAttivo] = (self._cursors[self._utenteAttivo][0],idpost)

        if(self._cursors[self._utenteAttivo][0] != posizione):
            self._spedisci('\C'+str(posizione)+'\\')
            #self._cursors[self._utenteAttivo] = (self._cursors[self._utenteAttivo][0]+len(dati),idpost)
        self._spedisci(dati)
        self._cursors[self._utenteAttivo] = (posizione+len(dati),idpost)
        print('posizione finale',self._cursors[self._utenteAttivo][0])
        
    def spedisci_rimozione(self,posizione,rimossi,idpost):
        '''
        Spedisce al server una segnalazione di rimozione testo, con puntatore e numero di
        caratteri che sono stati rimossi.
        '''
        if(self._cursors[self._utenteAttivo][1] != idpost):
            self._spedisci('\P'+str(idpost)+'\\')
            self._cursors[self._utenteAttivo] = (self._cursors[self._utenteAttivo][0],idpost)
        posizione += rimossi
        if(posizione != self._cursors[self._utenteAttivo][0]):
            self._spedisci('\C'+str(posizione)+'\\')
        self._spedisci('\D'+str(rimossi)+'\\')
        self._cursors[self._utenteAttivo] = (posizione,idpost)
        
class Receiver(QtCore.QThread):
    '''
    Classe che rappresenta il thread per ricevere i dati
    '''
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
            #print('ricevimento')
            try:
                buffer = self._socket.recv(1)
                try:
                    buffer = buffer.decode("utf-8")
                except:
                    buffer += self._socket.recv(1)
                    buffer = buffer.decode("utf-8")
                #print(buffer)
                self._coda.put(buffer)
            except:
                print('eccezione',sys.exc_info())
                self._stop = True
        self._stop = False
    
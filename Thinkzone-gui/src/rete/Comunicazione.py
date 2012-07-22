'''
Classe comunicatore che riesce a ricevere e inviare dati.
Ha dentro definite le meccaniche di comunicazione.
@author: stengun
'''
import queue,socket,sys,logging
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
    _registered = False
    _utenteAttivo = 0
    _userID = None
    _receive_thread = None
    _lastResponse = None
    _activePost = None
    #robe interne
    _barrier = None
    _cursors = {}
    def __init__(self):
        QtCore.QThread.__init__(self)
        self._socket = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        self._messaggi = queue.Queue(255)
        self._stop = False
        self._logger = logging.getLogger()
        _rimozione = QtCore.pyqtSignal(int,int,name='rimozione')
        _aggiunta = QtCore.pyqtSignal(int,str,name='aggiunta')
        _nuovopost = QtCore.pyqtSignal(int,name='nuovoPost')
        _selpost = QtCore.pyqtSignal(int,name='selectPost')
        _nuovutente = QtCore.pyqtSignal(int,int,name='cambiaUtente')
        self.blink_cursor = 0
        
    def run(self):
        self._cursors[0] = (0,0)
        while(not(self._stop)):
            try:
                messaggio = self._messaggi.get(True, None)
            except:
                self._logger.error("Errore nel prelievo del messaggio dal thread.",exc_info=True)
                #print('errore connettore',sys.exc_info())
                self._stop = True
                self._receive_thread.setTerminationEnabled()
                self._receive_thread._stop = True
                continue
            if(messaggio == ""):
                #self.disconnetti()
                self._receive_thread.setTerminationEnabled()
                self._receive_thread._stop = True
                self._stop = True
                continue
            #print('messaggio',messaggio)
            if(self._parseinput(messaggio) and self._utenteAttivo != self._userID):
                self._logger.debug('Utente %s, post %s. Scrivo %s in posizione %s',
                      self._utenteAttivo,
                      self._cursors[self._utenteAttivo][1],
                      messaggio,
                      self._cursors[self._utenteAttivo][0])
                
                self.emit(QtCore.SIGNAL('aggiunta(int,QString)'),self._cursors[self._utenteAttivo][0],messaggio)
                self._cursors[self._utenteAttivo] = (self._cursors[self._utenteAttivo][0]+1,self._cursors[self._utenteAttivo][1])
                
        try:
            self.disconnetti()
        except:
            self._logger.warning("Errore in fase di disconnessione: %s",sys.exc_info())
        finally:
            self._stop = False
           
    def _controller(self,controllo,messaggio):
        '''
        Questo metodo decide se il carattere "controllo" è un carattere di fine comando.
        Se questo test fallisce, imposta il valore di messaggio a None.
        '''
        if(controllo != '\\'):
            self._logger.critical("Errore nel formato dei messaggi!")
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
                [self.blink_cursor, controllo] = self._recvInt()
                if(self._userID != self._utenteAttivo):
                    print('cursore ',self.blink_cursor)
                    self._cursors[self._utenteAttivo] = (self.blink_cursor,self._activePost)
                    self._logger.debug("Utente %s.\n \
                                        Spostamento cursore all'indice %s",
                                        str(self._utenteAttivo),str(self.blink_cursor))
                messaggio = self._controller(controllo, messaggio)
                return False
            
            if(messaggio == 'D'): # Eliminazione
                print('caso D')
                [quantita, controllo] = self._recvInt()
                self._cursors[self._utenteAttivo] = ((self._cursors[self._utenteAttivo][0] - quantita),self._activePost)
                messaggio = self._controller(controllo, messaggio)
                if(messaggio == None):
                    return False
                else:
                    if(self._utenteAttivo != self._userID):
                        self.emit(QtCore.SIGNAL('rimozione(int,int)'),self._cursors[self._utenteAttivo][0],quantita)
                        self._logger.debug("Utente %s.\n \
                                            Rimossi %s caratteri.",
                                            str(self._utenteAttivo),str(quantita))
                return False
            
            if(messaggio== 'U'): # Selezione utente
                #print('caso U')
#                prece = self._utenteAttivo
                [self._utenteAttivo,controllo] = self._recvInt()
                try:
                    self._cursors[self._utenteAttivo]
                except:
                    self._logger.info("Registrato utente %s.",str(self._utenteAttivo))
                    self._cursors[self._utenteAttivo] = (0,0)
                self._logger.debug("L'utente %s è ora attivo.",str(self._utenteAttivo))
                #if(self._utenteAttivo != self._userID):
                #    self.cursore_locale = self.blink_cursor
                #if(prece != None):
                #    self.emit(QtCore.SIGNAL('cambiaUtente(int,int)'),prece,self._activePost)
                #    self._barrier.wait()
                messaggio = self._controller(controllo, messaggio)
                return False
            
            if(messaggio =='R'): #Risposta a una azione
                [self._lastResponse,controllo] = self._recvInt()
                self._logger.info("Risposta dal server: %s",str(self._lastResponse))
                messaggio = self._controller(controllo, messaggio)
                return False
            
            if(messaggio == 'P'):#selezione post
                #print('caso POST')
                [idpost,controllo] = self._recvInt()
                self._logger.debug("Utente %s, seleziona il post %s",str(self._utenteAttivo),str(idpost))
                messaggio = self._controller(controllo, messaggio)
                self._activePost = idpost
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
        if(self._registered):
            self._logger.info("Registrazione già effettuata.")
            return self._registered
        
        self._socket.connect((hostname,porta))
        self._receive_thread = Receiver(self._messaggi,self._socket)
        self._receive_thread.start()
        self._spedisci('\L1\\')
        self._spedisci(nickname+'\\'+password+'\\')
        messaggio = self._messaggi.get(True, None)
        if(self._parseinput(messaggio)):
            risposta = self._parseResponse(self._lastResponse)
            if(risposta[0]):
                print(risposta[1])
                return self._registered
        
        self._logger.info("Registrazione completata correttamente.")
        self._registered = True
        self.disconnetti()
        return self._registered
    
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
            risposta = self._parseResponse(self._lastResponse)
            if(risposta[0]):
                print(risposta[1])
                return False
        
        controllo = '\00'
        risposta = self._parseResponse(self._lastResponse)
        if(risposta[0]):
            errore = risposta[1]
            try:
                self.disconnetti()
            except:
                self._logger.error(errore,exc_info=True)
            finally:
                #print(errore)
                self._stop = True
                return False
        try:
            [self._userID,controllo] = self._recvInt()
        except:
            self._logger.critical("Connessione chiusa dal server!",exc_info=True)
            #print('connessione chiusa!',file=sys.stderr)
            sys.exit()
        
        #print('il tuo ID',self._userID,' controllo:',controllo)
        self._cursors[self._userID] = (0,0)
        self._logger.info("Connessione effettuata correttamente.\n \
                            ID utente %s",str(self._userID))
        return True
        
    def _parseResponse(self,response):
        '''
        Metodo che analizza i codici di risposta, e ritorna True se è un codice di errore.
        '''
        if(response == 0):
            self._logger.info("Response OK")
            return False,"OK"
        if(response == 1):
            self._logger.error("Impossibile effettuare il login, utente non registrato.")
            return True,"Utente non registrato."
        if(response == 2):
            self._logger.info("Utente già connesso al server.")
            self._registered = True
            return False,"Utente già connesso altrove"
        if(response == 3):
            self._logger.error("Impossibile effettuare il login, password errata.")
            return True,"Password Errata."
        self._logger.error("C'è stato un errore non catalogato. ID %s",str(response))
        return True,"Errore generico."
    
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
        #self._spedisci('\00')
#        print('posizione attuale',self._cursors[self._userID][0])
#        print('posizione da cui partire',posizione)
#        print('posizione da aggiungere',len(dati))
#        
        if(self._cursors[self._userID][1] != idpost):
            self._spedisci('\P'+str(idpost)+'\\')
            #self._cursors[self._utenteAttivo] = (self._cursors[self._utenteAttivo][0],idpost)

        if(self._cursors[self._userID][0] != posizione):
            self._spedisci('\C'+str(posizione)+'\\')
            #self._cursors[self._utenteAttivo] = (self._cursors[self._utenteAttivo][0]+len(dati),idpost)
        self._spedisci(dati)
        self._cursors[self._userID] = (posizione+len(dati),idpost)
        self._logger.debug("POST %s: Spedita aggiunta da %s di %s caratteri.",str(idpost),str(posizione),str(len(dati)))
        #print('posizione finale',self._cursors[self._userID][0])
        
    def spedisci_rimozione(self,posizione,rimossi,idpost):
        '''
        Spedisce al server una segnalazione di rimozione testo, con puntatore e numero di
        caratteri che sono stati rimossi.
        '''
        if(self._cursors[self._userID][1] != idpost):
            self._spedisci('\P'+str(idpost)+'\\')
            self._cursors[self._userID] = (self._cursors[self._userID][0],idpost)
        posizione += rimossi
        if(posizione != self._cursors[self._userID][0]):
            self._spedisci('\C'+str(posizione)+'\\')
        
        self._spedisci('\D'+str(rimossi)+'\\')
        self._cursors[self._userID] = (posizione-1,idpost)
        self._logger.debug("POST %s: Spedita rimozione di %s caratteri da %s",str(idpost),str(rimossi),str(posizione))
        
class Receiver(QtCore.QThread):
    '''
    Classe che rappresenta il thread per ricevere i dati
    '''
    _coda = queue.Queue()
    _stop = None
    _socket = None
    _logger = logging.getLogger("file_log")
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
                self._coda.put(buffer)
            except:
                self._logger.exception("Thread ricezione: chiusura thread in corso.")
                #print('eccezione in ricezione',sys.exc_info())
                self._stop = True
                self._coda = None
        self._stop = False
    
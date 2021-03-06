'''
Modulo che permette di ricevere e inviare dati.
Ha dentro definite le meccaniche di comunicazione.

Tutti i comandi sono nel formato **\\XY\\**
Qui di seguito sono elencati tutti i comandi supportati dal server:

Ricezione/Invio:
 * **\\Px\\** - Comando che indica il post attivo per i comandi successivi. *x* è un intero, indica l'id del post.
 * **\\Cx\\** - Indica un aggiornamento per il cursore della scrittura/lettura. *x* è un intero, indica la posizione del cursore.
 * **\\Dx\\** - Indica la cancellazione di qualcosa. *x* è un intero che dice quanti caratteri sono stati cancellati.

Ricezione:
 * **\\Rx\\**   - Response dal server per un dato comando. *x* è un intero che rappresenta il codice risposta.
 * **\\Ux\\**   - Indica l'utente attivo per i comandi successivi. *x* è un intero, rappresenta l'ID utente.
 * **\\Kx\y\\** - Indica la creazione di un post. *x* rappresenta l'id del post "padre", *y* l'id del post appena creato.

Invio:
 * **\\Lx\\** - Indica al server il comando di login. *x* rappresenta il comando. *0* per login, *1* per registrazione.
 * **\\Kx\\** - Indica al server che voglio creare un post. *x* è il padre del post.
'''
import queue,socket,sys,logging
from PyQt4 import QtCore

class comunicatore(QtCore.QThread):
    '''
    Generico comunicatore che imposta anche i thread per la ricezione.
    Effettua il parsing dei dati ricevuti e dei responses, loggando anche i messaggi di errore.
    '''

    _socket = None
    
    _posizione = None
    _messaggi = None
    _stop = None
    
    _registered = False

    _receive_thread = None
    _lastResponse = None

    #robe interne
    _barrier = None
    #_cursors = {}
    _postPlexer = None
    
    def __init__(self):
        QtCore.QThread.__init__(self)
        self._socket = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        self._messaggi = queue.Queue(255)
        self._stop = False
        self._logger = logging.getLogger()
        _refreshcursor = QtCore.pyqtSignal(int,name='refreshCursor')
        _applydelete = QtCore.pyqtSignal(int, name='applyDelete')
        _applyaggiunta = QtCore.pyqtSignal(str, name='applyAggiunta')
        _nuovopost = QtCore.pyqtSignal(int,name='nuovoPost')
        _selpost = QtCore.pyqtSignal(int,name='refreshPost')
        _nuovutente = QtCore.pyqtSignal(int,name='selectUser')
        _myid = QtCore.pyqtSignal(int,name='myId')
        
    def run(self):
        while(not(self._stop)):
            try:
                messaggio = self._messaggi.get(True, None)
            except:
                self._logger.error("Errore nel prelievo del messaggio dal thread.",exc_info=True)
                self._stop = True
                self._receive_thread.setTerminationEnabled()
                self._receive_thread._stop = True
                continue
            if(messaggio == ""):
                self._receive_thread.setTerminationEnabled()
                self._receive_thread._stop = True
                self._stop = True
                continue
            if(self._parseinput(messaggio)):
                self.emit(QtCore.SIGNAL('applyAggiunta(QString)'),messaggio)
                #self._postPlexer.applyAggiunta(messaggio)        
        try:
            self.disconnetti()
        except:
            self._logger.warning("Errore in fase di disconnessione: %s",sys.exc_info())
        finally:
            self._stop = False
           
    def _controller(self,controllo,messaggio):
        '''
        Questo metodo decide se il carattere *controllo* è un carattere di fine comando.
        Se questo test fallisce, imposta il valore di messaggio a None.
        '''
        if(controllo != '\\'):
            self._logger.critical("Errore nel formato dei messaggi!")
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
        altrimenti ritorna False quando si tratta di un comando.
        Per ogni comando effettua l'azione corrispondente prima di uscire.
        '''
        controllo = '\00'
        
        if(messaggio == '\\'): # selezione comando
            messaggio = self._messaggi.get(True,None)
            if(messaggio == '\\'): #non è un comando.
                return True
            
            if(messaggio == 'C'): #Modifica cursore
                [cursore, controllo] = self._recvInt()
                self.emit(QtCore.SIGNAL('refreshCursor(int)'),cursore)
                messaggio = self._controller(controllo, messaggio)
                return False
            
            if(messaggio == 'D'): # Eliminazione
                [quantita, controllo] = self._recvInt()
                messaggio = self._controller(controllo, messaggio)
                if(messaggio == None):
                    return False
                else:
                    self.emit(QtCore.SIGNAL('applyDelete(int)'),quantita)
                return False
            
            if(messaggio== 'U'): # Selezione utente
                [utente,controllo] = self._recvInt()
                self.emit(QtCore.SIGNAL('selectUser(int)'),utente)
                messaggio = self._controller(controllo, messaggio)
                return False
            
            if(messaggio =='R'): #Risposta a una azione
                [self._lastResponse,controllo] = self._recvInt()
                self._logger.info("Risposta dal server: %s",str(self._lastResponse))
                messaggio = self._controller(controllo, messaggio)
                return False
            
            if(messaggio == 'P'):#selezione post
                [idpost,controllo] = self._recvInt()
                messaggio = self._controller(controllo, messaggio)
                self.emit(QtCore.SIGNAL('refreshPost(int)'),idpost)
                #self._postPlexer.refreshPost(idpost)
                return False
            
            if(messaggio == 'K'):
                [parent,controllo] = self._recvInt()
                [idpost,controllo] = self._recvInt()
                self.emit(QtCore.SIGNAL('nuovoPost(int)'),idpost)
                self._barrier.wait()
                return False
            
        return True
    
    def registrati(self,hostname,porta,nickname,password):
        '''
        Metodo per effettuare una registrazione ad un server.
        Le situazioni di errore o di normalità vengono scritte nel file di log.
        '''
        if(self._registered):
            self._logger.info("Registrazione già effettuata.")
            return self._registered
        
        try:
            self._socket.connect((hostname,porta))
        except:
            print('Errore registrazione: ',sys.exc_info(),sys.stderr)
            return self._registered
        
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
        try:
            self._socket.connect((hostname,porta))
        except:
            self._stop = True
            self.disconnetti()
            return False
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
            [_userID,controllo] = self._recvInt()
            if(self._controller(controllo, messaggio) == None):
                raise BaseException
        except:
            self._logger.critical("Connessione chiusa dal server!",exc_info=True)
            sys.exit()
        
        self.emit(QtCore.SIGNAL('myId(int)'),_userID)
        self._logger.info("Connessione effettuata correttamente.\
                            ID utente %s",str(_userID))
        return True
        
    def _parseResponse(self,response):
        '''
        Metodo che analizza i codici di risposta, e ritorna True se è un codice di errore.
        Codici di risposta:
            0 - Tutto ok
            1 - (login) Errore: utente non registrato
            2 - (login) Errore: utente già connesso
            3 - (login) Errore: password sbagliata
        
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
        try:
            self._socket.shutdown(socket.SHUT_RDWR)
            self._socket.close()
        except:
            self._logger.warning("Errore in disconnessione, resetto il socket.")
        self._socket = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        if(self._receive_thread != None):
            self._receive_thread.setTerminationEnabled(True)
            self._receive_thread._stop = True
        self._postPlexer._userID = None
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
        if(self._postPlexer.myActivePost() != idpost):
            self._spedisci('\P'+str(idpost)+'\\')
            self._postPlexer.refreshmyPost(idpost)

        if(self._postPlexer.myActiveCursor() != posizione):
            self._spedisci('\C'+str(posizione)+'\\')

        self._spedisci(dati)
        self._postPlexer.refreshmyCursor(posizione+len(dati))
        self._logger.debug("POST %s: Spedita aggiunta da %s di %s caratteri.",str(idpost),str(posizione),str(len(dati)))
        
    def spedisci_rimozione(self,posizione,rimossi,idpost):
        '''
        Spedisce al server una segnalazione di rimozione testo, con puntatore e numero di
        caratteri che sono stati rimossi.
        '''
        if(self._postPlexer.myActivePost() != idpost):
            self._spedisci('\P'+str(idpost)+'\\')
            self._postPlexer.refreshmyPost(idpost)
        posizione += rimossi
        if(posizione != self._postPlexer.myActiveCursor()):
            self._spedisci('\C'+str(posizione)+'\\')
        
        self._spedisci('\D'+str(rimossi)+'\\')
        self._postPlexer.refreshmyCursor(posizione)
        self._logger.debug("POST %s: Spedita rimozione di %s caratteri da %s",str(idpost),str(rimossi),str(posizione))
        
class Receiver(QtCore.QThread):
    '''
    Classe che rappresenta il thread per ricevere i dati
    L'unica cosa che fa è ricevere dei dati da un socket e inserirli dentro una coda.
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
                self._stop = True
                self._coda = None
        self._stop = False
    
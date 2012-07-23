'''
Crea la finestra principale del programma Thinkzone.
In questa finestra è possibile creare e modificare conversazioni e post.
Questo modulo contiene soltanto la classe da istanziare per creare la finestra principale.
'''
import logging
from gui import finestraprincipale,loginDialog,aboutDialog
from PyQt4 import QtGui, QtCore
from utils import PostWidget
from rete import Comunicazione
from threading import Barrier

class mainwindow(QtGui.QMainWindow,finestraprincipale.Ui_MainWindow):
    '''
    Classe che crea una finestra principale per il programma.
    '''
    #dialogs e finestre
    _loginDialog = None
    _aboutDialog = None
    
    _postPlexer = None
    #thread e altro
    _logger = None
    _barrier = None
    
    _connettore = Comunicazione.comunicatore()
    __VERSION__ = None

    def __init__(self,version,parent = None):
        '''
        Costruisce una finestra principale di Thinkzone.
        Ha bisogno di alcuni parametri per avviarsi.
        @param version: La versione del programma. 
        '''
        self.__VERSION__ = version
        logging.basicConfig(filename="thinkzone_gui.log",format='%(asctime)s | Loglevel: %(levelname)s | %(message)s', datefmt='%m/%d/%Y %I:%M:%S %p')
        self._logger = logging.getLogger()
        self._logger.setLevel(logging.DEBUG)
        self._logger.info("Inizio nuova sessione")
        QtGui.QMainWindow.__init__(self,parent)
        self.ui = finestraprincipale.Ui_MainWindow()
        self.setupUi(self)
        spacerItem = QtGui.QSpacerItem(20, 40, QtGui.QSizePolicy.Minimum, QtGui.QSizePolicy.Expanding)
        self.layoutTextarea.addItem(spacerItem)
        
        #setup finestre di dialogo
        self._loginDialog = loginDialog.Login(self)
        self._aboutDialog = aboutDialog.aboutDial(self)
        self._aboutDialog.labelVersion.setText("Version: "+self.__VERSION__)
        #connettori
        self._barrier = Barrier(2, timeout=200)
        self._connettore._barrier = self._barrier
        self._postPlexer = PostPlexer(self._connettore)
        self._connettore._postPlexer = self._postPlexer
        #connessione di tutti i segnali
        QtCore.QObject.connect(self.buttonCrea,QtCore.SIGNAL('pressed()'), self._inviaCreazione)
        QtCore.QObject.connect(self.actionLogin,QtCore.SIGNAL("triggered()"),self._loginDialog.show)
        QtCore.QObject.connect(self.actionInformazioni_su,QtCore.SIGNAL("triggered()"),self._aboutDialog.show)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('nuovoPost(int)'),self._creapost,2)
        self._creapost(0)
        self.layout_titolo.addWidget(self._postPlexer._postids[0])
    
    def _inviaCreazione(self):
        '''
        Crea nella finestra un nuovo post.
        '''
        atti = self._postPlexer.myActivePost()
        if(atti == None):
            atti = 0
        self._logger.debug("Creazione nuovo post.")
        self._connettore._spedisci('\K'+str(atti)+'\\') # WARNING da ricontrollare meglio.

        
    def _creapost(self,idpost):
        '''
        Parsing degli ID dei post. Crea nuovi post, seleziona post precedenti e dice quando non esistono.
        '''          
        self._logger.debug("Creazione di un nuovo post con ID: "+str(idpost))
        textArea = PostWidget.postWidget(idpost)
        QtCore.QObject.connect(textArea,QtCore.SIGNAL('testoRimosso(int,int,int)'), self._connettore.spedisci_rimozione)
        QtCore.QObject.connect(textArea,QtCore.SIGNAL('testoAggiunto(int,QString,int)'), self._connettore.spedisci_aggiunta)
        self._postPlexer._postids[idpost] = textArea
        if(idpost != 0):
            self.layoutTextarea.addWidget(textArea)
            self._barrier.wait()


class PostPlexer(QtCore.QObject):
    '''
    Comoda classe per il multiplexing dei post.
    Offre una interfaccia unica per cambiare utenti, selezionare i post
    e ottenere i post attivi. Contiene anche l'id dell'utente locale.
    '''
    _connettore = None
    _userID = None
    _utenteAttivo = None
    _utentePrecedente = None
    _postids = {}
    _cursors = {}
    _logger = None
    def __init__(self,connettore,parent = None):
        QtCore.QObject.__init__(self)
        self._connettore = connettore
        self._logger = logging.getLogger()
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('selectUser(int)'),self.selectUser,2)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('refreshPost(int)'),self.refreshPost,2)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('refreshCursor(int)'),self.refreshCursor,2)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('applyDelete(int)'),self.applyDelete,2)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('applyAggiunta(QString)'),self.applyAggiunta,2)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('myId(int)'),self._myId,2)
        
    def _myId(self,myid):
        self._userID = myid
        self.addUser(myid)
        
    def _select(self,postid):
        userdata = (self._getCursor(),postid)
        self._cursors[self._utenteAttivo] = userdata
        self._logger.debug("Selezionato il post"+str(postid))
    
    def _getCursor(self):
        return self._cursors[self._utenteAttivo][0]
    
    def _postConnect(self):
        self._select(self.postAttuale())
        self._logger.debug("Utente %s, seleziona il post %s",str(self._utenteAttivo),str(self.postAttuale()))
    
    def refreshCursor(self,cursor):
        '''
        Aggiorna il cursore dell'utente attivo.
        '''
        userdata = (cursor,self.postAttuale())
        self._cursors[self._utenteAttivo] = userdata
        self._logger.debug("Utente %s. \
                            Spostamento cursore all'indice %s",
                            str(self._utenteAttivo),str(self._getCursor()))
    
    def refreshmyCursor(self,cursor):
        '''
        Aggiorna il cursore per l'utente locale.
        '''
        userdata = (cursor,self.myActivePost())
        self._cursors[self._userID] = userdata
        
    def refreshPost(self,postid):
        '''
        aggiorna il post dell'utente attivo.
        '''
        userdata = (self._getCursor(),postid)
        self._cursors[self._utenteAttivo] = userdata
        self._postConnect()
    
    def refreshmyPost(self,postid):
        '''
        Modifica l'id del post attivo per l'utente locale.
        '''
        userdata = (self.myActiveCursor(),postid)
        self._cursors[self._userID] = userdata
        
    def applyDelete(self,quantita):
        '''
        Cancella i caratteri dal post selezionato.
        Se l'utente attivo è quello locale questo metodo non fa nulla.
        '''
        if(self._userID != self._utenteAttivo):
            cursore = self._getCursor()
            post = self._postids[self.postAttuale()]
            cursore = cursore - quantita
            post.rimuoviTesto(cursore,quantita)
            self.refreshCursor(cursore)
            self._logger.debug("Utente %s. Rimossi %s caratteri.",
                               str(self._utenteAttivo),str(quantita))
    
    def applyAggiunta(self,stringa):
        '''
        Aggiunge i caratteri ai post selezionati.
        Se l'utente attivo è quello locale, questo metodo non fa nulla.
        '''
        if(self._userID != self._utenteAttivo):
            cursore = self._getCursor()
            post = self._postids[self.postAttuale()]
            post.aggiungiTesto(cursore,stringa)
            cursore = cursore + len(stringa)
            self.refreshCursor(cursore)
            self._logger.debug('Utente %s, post %s. Scrivo %s. Posizione finale %s',
                          self._utenteAttivo,
                          self.postAttuale(),
                          stringa,
                          self._getCursor())
        
    def selectUser(self,utente):
        '''
        Seleziona un utente per marcarlo come attivo.
        '''
        try:
            self._cursors[utente]
        except:
            self.addUser(utente)
        finally:
            self._utentePrecedente = self._utenteAttivo
            self._utenteAttivo = utente
        
        self._logger.debug("L'utente %s è ora attivo.",str(self._utenteAttivo))
        self._postConnect()
    
    def addUser(self,utente):
        '''
        Aggiunge un nuovo utente alla struttura dati.
        '''
        self._cursors[utente] = (0,0)
    
    def postAttuale(self):
        '''
        Restituisce l'indice del post dell'utente attivo.
        '''
        return self._cursors[self._utenteAttivo][1]
    
    def myActivePost(self):
        '''
        Restituisce l'indice del post attivo per l'utente locale.
        '''
        return self._cursors[self._userID][1]
    
    def myActiveCursor(self):
        '''
        Restituisce il cursore dell'utente locale.
        '''
        return self._cursors[self._userID][0]
        
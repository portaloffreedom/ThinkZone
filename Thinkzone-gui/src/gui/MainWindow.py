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
    #widget personalizzati
    _risposte = []
    _postids = {}
    #thread e altro
    _logger = None
    _barrier = None
    _settato = False
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

        #connessione di tutti i segnali
        QtCore.QObject.connect(self.buttonCrea,QtCore.SIGNAL('pressed()'), self._inviaCreazione)
        QtCore.QObject.connect(self.actionLogin,QtCore.SIGNAL("triggered()"),self._loginDialog.show)
        QtCore.QObject.connect(self.actionInformazioni_su,QtCore.SIGNAL("triggered()"),self._aboutDialog.show)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('nuovoPost(int)'),self._creapost,2)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('selectPost(int)'),self._selpost,2)
        #QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('cambiaUtente(int,int)'),self._cambiautente,2)
        self._creapost(0)
        self.layout_titolo.addWidget(self._postids[0])
        self._barrier = Barrier(2, timeout=200)
        self._connettore._barrier = self._barrier

    
    def _inviaCreazione(self):
        '''
        Crea nella finestra un nuovo post.
        '''
        atti = self._connettore._cursors[self._connettore._utenteAttivo][1]
        if(atti == None):
            atti = 0
#        if(atti == 0):
#            if(self.titoloEdit.isEnabled()):
#                self._logger.debug("Spedizione del titolo: "+self.postids[0])
#                self._creapost(0)
#                self._connettore._spedisci('\P0\\'+self.titoloEdit.text())
#                self.titoloLabel.setText(self.titoloEdit.text())
#                self.titoloEdit.setText('')
#                self.titoloEdit.setDisabled(True)
#                self.titoloEdit.setHidden(True)
        self._logger.debug("Creazione nuovo post.")
        #self._connettore._spedisci('\K0\\') #warning
        self._connettore._spedisci('\K'+str(atti)+'\\') # WARNING da ricontrollare meglio.
    
    def _selectPost(self,idpost):
        selezionato = self._postids[idpost]
        if(not(selezionato._selected)):
            selezionato._selected = True
            QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('rimozione(int,int)'),selezionato.rimuoviTesto,2)
            QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),selezionato.aggiungiTesto,2)
            self._logger.debug("Selezionato il post"+str(idpost))
    
    def _deselectPost(self,idpost):
        #if(idpost == 0):
        #    QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),self._setTitolo)
        #    return
        selezionato = self._postids[idpost]
        if(selezionato._selected):
            selezionato._selected = False
            QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('rimozione(int,int)'),selezionato.rimuoviTesto)
            QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),selezionato.aggiungiTesto)
            self._logger.debug("Deselezionato il post"+str(idpost))
    
#    def _cambiautente(self,precedente,attuale):
#        #print('CAMBIO UTENTE')
#        if(precedente == 0):
#            #self.titoloLabel.setText('Titolo: '+self.titoloLabel.text())
#            QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),self._setTitolo)
#            #self._connettore._activePost = None
#        self._barrier.wait()
    
    def _selpost(self,idpost):
        precedente = self._connettore._cursors[self._connettore._utenteAttivo][1]
        if(precedente != None):
            self._deselectPost(precedente)
        self._selectPost(idpost)
#        else:
#            if(not(self._settato)):
#                QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),self._setTitolo,2)
#                self._settato = True
        self._connettore._cursors[self._connettore._utenteAttivo] = (self._connettore._cursors[self._connettore._utenteAttivo][0],idpost)
        
        self._barrier.wait()
    
#    def _setTitolo(self,indi,stringa):
#        if(stringa != ''):
#            self._logger.debug("Impostazione titolo discussione: "+stringa)
#            self.titoloLabel.setText(self.titoloLabel.text()+stringa)
#            self.titoloEdit.setDisabled(True)
#            self.titoloEdit.setHidden(True)
        
    def _creapost(self,idpost):
        '''
        Parsing degli ID dei post. Crea nuovi post, seleziona post precedenti e dice quando non esistono.
        '''          
        self._logger.debug("Creazione di un nuovo post con ID: "+str(idpost))
        textArea = PostWidget.postWidget(idpost)
        self._postids[idpost] = textArea
        self.layoutTextarea.addWidget(textArea)
        QtCore.QObject.connect(textArea,QtCore.SIGNAL('testoRimosso(int,int,int)'), self._connettore.spedisci_rimozione)
        QtCore.QObject.connect(textArea,QtCore.SIGNAL('testoAggiunto(int,QString,int)'), self._connettore.spedisci_aggiunta)
        if(idpost != 0):
            self._barrier.wait()

        
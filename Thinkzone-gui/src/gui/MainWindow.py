'''
Crea e avvia il programma principale Thinkzone.
@author: stengun
'''
import sys
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
    _barrier = None
    _settato = False
    _connettore = Comunicazione.comunicatore()
    __VERSION__ = "0.0.7"

    def __init__(self,parent = None):
        QtGui.QMainWindow.__init__(self,parent)
        #setup interfaccia principale
        self.ui = finestraprincipale.Ui_MainWindow()
        self.setupUi(self)
        #self.layoutTextarea.addWidget(self._textArea)
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
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('cambiaUtente(int,int)'),self._cambiautente,2)
        self._barrier = Barrier(2, timeout=200)
        self._connettore._barrier = self._barrier
    
    def _inviaCreazione(self):
        '''
        Crea nella finestra un nuovo post.
        '''
        atti = self._connettore._cursors[self._connettore._utenteAttivo][1]
        if(atti == None):
            atti = 0
        if(atti == 0):
            self._connettore._spedisci('\P0\\'+self.titoloEdit.text())
            self.titoloLabel.setText(self.titoloEdit.text())
            self.titoloEdit.setDisabled(True)
            
        self._connettore._spedisci('\K'+str(atti)+'\\') #FIXME ci sarà un bel bug qui dentro
    
    def _selectPost(self,idpost):
        selezionato = self._postids[idpost]
        if(not(selezionato._selected)):
            selezionato._selected = True
            QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('rimozione(int,int)'),selezionato.rimuoviTesto,2)
            QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),selezionato.aggiungiTesto,2)
            print('selezionato il post '+str(idpost))
    
    def _deselectPost(self,idpost):
        if(idpost == 0):
            QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),self._setTitolo)
            return
        selezionato = self._postids[idpost]
        if(selezionato._selected):
            selezionato._selected = False
            QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('rimozione(int,int)'),selezionato.rimuoviTesto)
            QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),selezionato.aggiungiTesto)
            print('deselezionato il post '+str(idpost))
    
    def _cambiautente(self,precedente,attuale):
        print('CAMBIO UTENTE')
        if(precedente == 0):
            #self.titoloLabel.setText('Titolo: '+self.titoloLabel.text())
            QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),self._setTitolo)
            #self._connettore._activePost = None
        self._barrier.wait()
    
    def _selpost(self,idpost):
        #if(self._connettore._cursors[self._connettore._utenteAttivo][1] != idpost):
        if(idpost !=0):
            precedente = self._connettore._cursors[self._connettore._utenteAttivo][1]
            if(precedente != None):
                self._deselectPost(precedente)
            self._selectPost(idpost)
        else:
            if(not(self._settato)):
                print('connesso titolo')
                QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),self._setTitolo,2)
                self._settato = True
        self._connettore._cursors[self._connettore._utenteAttivo] = (self._connettore._cursors[self._connettore._utenteAttivo][0],idpost)
        
        self._barrier.wait()
    
    def _setTitolo(self,indi,stringa):
        print('arrivato titolo: ',stringa)
        self.titoloLabel.setText(self.titoloLabel.text()+stringa)
        
    def _creapost(self,idpost):
        '''
        Parsing degli ID dei post. Crea nuovi post, seleziona post precedenti e dice quando non esistono.
        '''          
        if(self._connettore._cursors[self._connettore._utenteAttivo][1] == 0):
            print('CREO POST E TOLGO IL TITOLO')
            QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),self._setTitolo)
        if(idpost == 0):
            #return
            idpost +=1
            #TODO inserire un elemento in lista conversazione
            #pass
        print('creazione nuovo post con ID',idpost)
        textArea = PostWidget.postWidget(idpost)
        self._postids[idpost] = textArea
        self.layoutTextarea.addWidget(textArea)
        QtCore.QObject.connect(textArea,QtCore.SIGNAL('testoRimosso(int,int,int)'), self._connettore.spedisci_rimozione)
        QtCore.QObject.connect(textArea,QtCore.SIGNAL('testoAggiunto(int,QString,int)'), self._connettore.spedisci_aggiunta)
        self._barrier.wait()
        #print('Post selezionato:',idpost)
        #self._connettore._activePost = idpost
        

if __name__ == '__main__':
    app = QtGui.QApplication(sys.argv)
    finestra = mainwindow()
    finestra.show()
    app.exec()
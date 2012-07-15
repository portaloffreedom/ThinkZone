'''
Created on 15/lug/2012

@author: stengun
'''
import sys
from gui import finestraprincipale,loginDialog,aboutDialog
from PyQt4 import QtGui, QtCore
from utils import PostWidget
from rete import Comunicazione
from random import Random

class mainwindow(QtGui.QMainWindow,finestraprincipale.Ui_MainWindow):
    '''
    classdocs
    '''
    #dialogs e finestre
    _loginDialog = None
    _aboutDialog = None
    #widget personalizzati
    _risposte = []
    _postids = {}
    #thread e altro
    _connettore = Comunicazione.comunicatore()
    VERSION = "0.0.6"

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
        self._aboutDialog.labelVersion.setText("Version: "+self.VERSION)
        
        #connessione di tutti i segnali
        QtCore.QObject.connect(self.buttonCrea,QtCore.SIGNAL('pressed()'), self._creapost)
        QtCore.QObject.connect(self.actionLogin,QtCore.SIGNAL("triggered()"),self._loginDialog.show)
        QtCore.QObject.connect(self.actionInformazioni_su,QtCore.SIGNAL("triggered()"),self._aboutDialog.show)
        #QtCore.QObject.connect(self._textArea,QtCore.SIGNAL('testoRimosso(int,int)'), self._connettore.spedisci_rimozione)
        #QtCore.QObject.connect(self._textArea,QtCore.SIGNAL('testoAggiunto(int,QString)'), self._connettore.spedisci_aggiunta)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('nuovoPost(int)'),self._parsepost,2)
        
        '''
        Constructor
        '''
    
    def _creapost(self):
        atti = self._connettore._activePost
        if(atti == None):
            atti = 0
        self._connettore._spedisci('\K'+str(atti)+'\\') #FIXME ci sarà un bel bug qui dentro
    
    def _parsepost(self,idpost):
        no = False
        try:
            precedente = self._postids[self._connettore._activePost]
        except:
            no = True
            print('Precedente nullo',sys.exc_info())
            
        if(not(no) and self._connettore._activePost != idpost):
            precedente = self._postids[self._connettore._activePost]
            QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('rimozione(int,int)'),precedente.rimuoviTesto,2)
            QtCore.QObject.disconnect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),precedente.aggiungiTesto,2)
            QtCore.QObject.disconnect(precedente,QtCore.SIGNAL('testoRimosso(int,int,int)'), self._connettore.spedisci_rimozione)
            QtCore.QObject.disconnect(precedente,QtCore.SIGNAL('testoAggiunto(int,QString,int)'), self._connettore.spedisci_aggiunta)
        
        try:
            textArea = self._postids[idpost]
        except:
            if(idpost == 0):
                idpost +=1
                #TODO inserire un elemento in lista conversazione
                pass
            print('creazione nuovo post con ID',idpost)
            textArea = PostWidget.postWidget(idpost)
            self._postids[idpost] = textArea
            self.layoutTextarea.addWidget(textArea)
        finally:
            QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('rimozione(int,int)'),textArea.rimuoviTesto,2)
            QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),textArea.aggiungiTesto,2)
            QtCore.QObject.connect(textArea,QtCore.SIGNAL('testoRimosso(int,int,int)'), self._connettore.spedisci_rimozione)
            QtCore.QObject.connect(textArea,QtCore.SIGNAL('testoAggiunto(int,QString,int)'), self._connettore.spedisci_aggiunta)
        

if __name__ == '__main__':
    app = QtGui.QApplication(sys.argv)
    finestra = mainwindow()
    finestra.show()
    app.exec()
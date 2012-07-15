'''
Created on 15/lug/2012

@author: stengun
'''
import sys
from gui import finestraprincipale,loginDialog,aboutDialog
from PyQt4 import QtGui, QtCore
from utils import PostArea
from rete import Comunicazione

class mainwindow(QtGui.QMainWindow,finestraprincipale.Ui_MainWindow):
    '''
    classdocs
    '''
    #dialogs e finestre
    _loginDialog = None
    _aboutDialog = None
    #widget personalizzati
    _textArea = None
    #thread e altro
    _connettore = Comunicazione.comunicatore()
    VERSION = "0.0.6"

    def __init__(self,parent = None):
        QtGui.QMainWindow.__init__(self,parent)
        #setup interfaccia principale
        self.ui = finestraprincipale.Ui_MainWindow()
        self.setupUi(self)
        self._textArea = PostArea.Post()
        self.layoutTextarea.addWidget(self._textArea)
        #setup finestre di dialogo
        self._loginDialog = loginDialog.Login(self)
        self._aboutDialog = aboutDialog.aboutDial(self)
        self._aboutDialog.labelVersion.setText("Version: "+self.VERSION)
        
        #connessione di tutti i segnali
        QtCore.QObject.connect(self.actionLogin,QtCore.SIGNAL("triggered()"),self._loginDialog.show)
        QtCore.QObject.connect(self.actionInformazioni_su,QtCore.SIGNAL("triggered()"),self._aboutDialog.show)
        QtCore.QObject.connect(self._textArea,QtCore.SIGNAL('testoRimosso(int,int)'), self._connettore.spedisci_rimozione)
        QtCore.QObject.connect(self._textArea,QtCore.SIGNAL('testoAggiunto(int,QString)'), self._connettore.spedisci_aggiunta)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('rimozione(int,int)'),self._textArea.rimuoviTesto,2)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),self._textArea.aggiungiTesto,2)
        '''
        Constructor
        '''
        

if __name__ == '__main__':
    app = QtGui.QApplication(sys.argv)
    finestra = mainwindow()
    finestra.show()
    app.exec()
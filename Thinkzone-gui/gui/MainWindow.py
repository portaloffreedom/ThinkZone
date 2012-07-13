'''
Created on 11/lug/2012
Finestra principale che si occupa di connettere le componenti della GUI e i widget personalizzati.
@author: stengun
'''
import sys
from utils import PostArea
from rete import Comunicazione
from gui import finestraprincipale,aboutDialog
from PyQt4 import QtGui, QtCore, Qt

class FinestraPrincipale(QtGui.QMainWindow, finestraprincipale.Ui_finestraprincipale):
    '''
    Classe per la finestra principale.
    '''
    _textextended = None
    _connettore = None
    _aboutwindow = None
    VERSION = "0.0.2"
    def __init__(self, parent = None):
        QtGui.QMainWindow.__init__(self, parent)
        self._textextended = PostArea.Post()
        self.ui = finestraprincipale.Ui_finestraprincipale()
        self._aboutwindow = aboutDialog.aboutDial()
        self.setupUi(self)
        self.serverBox.addItems(['Server personalizzato','localhost:4242','192.168.0.42:4242','portaloffreedom.is-a-geek.org:4242'])
        QtCore.QObject.connect(self.serverBox,QtCore.SIGNAL('currentIndexChanged(QString)'),self.cambioindici)
        self.scroll_layout_2.addWidget(self._textextended)
        self._connettore = Comunicazione.comunicatore()
        self._aboutwindow.labelVersion.setText("Version: "+self.VERSION)
        QtCore.QObject.connect(self.actionAbout_Thinkzone,QtCore.SIGNAL('triggered()'),self._aboutwindow.show),
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('rimozione(int,int)'),self._textextended.rimuoviTesto,2)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),self._textextended.aggiungiTesto,2)
        QtCore.QObject.connect(self.buttonConnect,QtCore.SIGNAL('pressed()'), self.connetti)
        QtCore.QObject.connect(self._textextended,QtCore.SIGNAL('testoRimosso(int,int)'), self.dati_rimossi)
        QtCore.QObject.connect(self._textextended,QtCore.SIGNAL('testoAggiunto(int,QString)'), self.dati_aggiunti)
        QtCore.QObject.connect(self.usernameEdit, QtCore.SIGNAL('textEdited(QString)'),self._abilitaLogin)
        QtCore.QObject.connect(self.passwordEdit, QtCore.SIGNAL('textEdited(QString)'),self._abilitaLogin)
        QtCore.QObject.connect(self.buttonRegister,QtCore.SIGNAL('pressed()'), self.registrati)
        
    def _abilitaLogin(self):
        pwdtext = self.passwordEdit.text()
        usertext = self.usernameEdit.text()
        if(pwdtext != '' and usertext != ''):
            self.buttonConnect.setEnabled(True)
            self.buttonRegister.setEnabled(True)
        else:
            self.buttonConnect.setEnabled(False)
            self.buttonRegister.setEnabled(False)
    
    def cambioindici(self,elemento):
        indes = elemento.find(':')
        host= elemento[:indes]
        porta = elemento[indes+1:]
        print(host,porta)
        if(elemento == 'Server personalizzato'):
            self.widget_hostname.setEnabled(True)
        else:
            self.hostEdit.setText(host)
            self.portaEdit.setText(porta)
            self.widget_hostname.setEnabled(False)
    
    def dati_rimossi(self,posizione,rimossi):
        #print('ci sono')
        self._connettore.spedisci_rimozione(posizione,rimossi)
    
    def dati_aggiunti(self,posizione,aggiunti):
        self._connettore.spedisci_aggiunta(posizione,aggiunti)
    
    def registrati(self):
        hostname = self.hostEdit.text()
        nickname = self.usernameEdit.text()
        password = self.passwordEdit.text()
        porta = self.portaEdit.text()
        porta = porta.encode()
        if(porta == '' or hostname == ''):
            print('Non puoi avere un campo vuoto su Host e Porta!',file=sys.stderr)
            return
        porta = int(porta)
        self._connettore.registrati(hostname, porta, nickname, password)
    
    def connetti(self):
        porta = self.portaEdit.text()
        porta = porta.encode()
        hostname = self.hostEdit.text()
        if(porta == '' or hostname == ''):
            print('Non puoi avere un campo vuoto su Host e Porta!',file=sys.stderr)
            return
        porta = int(porta)
        self._connettore.connetti(hostname, porta,self.usernameEdit.text(),self.passwordEdit.text())
        self._connettore.start()

        
if __name__ == '__main__':
    app = QtGui.QApplication(sys.argv)
    finestra = FinestraPrincipale()
    finestra.show()
    app.exec()
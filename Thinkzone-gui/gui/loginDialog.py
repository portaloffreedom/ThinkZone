'''
Dialog window per il login. Viene richiamato dalla finestra principale.
Costruisce la finestra di login dalla classe di PyQt4
@author: stengun
'''
import sys
from utils import PostArea
from gui import login,aboutDialog
from PyQt4 import QtGui, QtCore

class Login(QtGui.QDialog, login.Ui_Dialog):
    '''
    Classe per la finestra principale.
    Eredita da login.Ui_Dialog. Costruisce e connette tutti i componenti della finestra
    di login. Modificare questo file se si vuole aggiungere nuovi connettori o widget
    particolari.
    '''
    _connettore = None
    _parent = None
    
    def __init__(self, parent = None):
        QtGui.QDialog.__init__(self, parent)
        self._parent = parent
        self._textextended = PostArea.Post()
        self._aboutwindow = aboutDialog.aboutDial()
        self._connettore = parent._connettore
        self.setupUi(self)
        self.serverBox.addItems(['Server personalizzato','localhost:4242','192.168.0.42:4242','portaloffreedom.is-a-geek.org:4242'])

        QtCore.QObject.connect(self.serverBox,QtCore.SIGNAL('currentIndexChanged(QString)'),self.cambioindici)
        QtCore.QObject.connect(self.buttonConnect,QtCore.SIGNAL('pressed()'), self.connetti)
        QtCore.QObject.connect(self.usernameEdit, QtCore.SIGNAL('textEdited(QString)'),self._abilitaLogin)
        QtCore.QObject.connect(self.passwordEdit, QtCore.SIGNAL('textEdited(QString)'),self._abilitaLogin)
        QtCore.QObject.connect(self.buttonRegister,QtCore.SIGNAL('pressed()'), self.registrati)
        
    def _abilitaLogin(self):
        '''
        Imposta attivato o disattivato i pulsanti per connettersi (o registrarsi)
        '''
        pwdtext = self.passwordEdit.text()
        usertext = self.usernameEdit.text()
        if(pwdtext != '' and usertext != ''):
            self.buttonConnect.setEnabled(True)
            self.buttonRegister.setEnabled(True)
        else:
            self.buttonConnect.setEnabled(False)
            self.buttonRegister.setEnabled(False)
    
    def cambioindici(self,elemento):
        '''
        Viene chiamata quando lo spinner per la selezione dei server cambia la sua selezione.
        '''
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
    
#    def dati_rimossi(self,posizione,rimossi):
#        #print('ci sono')
#        self._connettore.spedisci_rimozione(posizione,rimossi)
#    
#    def dati_aggiunti(self,posizione,aggiunti):
#        self._connettore.spedisci_aggiunta(posizione,aggiunti)
    
    def registrati(self):
        '''
        Metodo chiamato quando si preme il pulsante di registrazione a un server.
        Invia nome utente e password inseriti come dati di registrazione.
        '''
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
        '''
        Metodo chiamato quando si preme il pulsante di connessione.
        '''
        porta = self.portaEdit.text()
        porta = porta.encode()
        hostname = self.hostEdit.text()
        if(porta == '' or hostname == ''):
            print('Non puoi avere un campo vuoto su Host e Porta!',file=sys.stderr)
            return
        porta = int(porta)
        self._connettore.connetti(hostname, porta,self.usernameEdit.text(),self.passwordEdit.text())
        self._connettore.start()
        self._parent._connettore = self._connettore
        self.close()

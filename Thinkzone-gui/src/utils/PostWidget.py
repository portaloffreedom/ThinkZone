'''
Created on 15/lug/2012
Widget che rappresenta un post per il programma.
I post hanno un ID che li rende unici.
Incorpora una textarea modificata per consentire a più utenti di accedervi contemporaneamente.
@author: stengun
'''
from utils import PostArea
from PyQt4 import QtGui, QtCore

class postWidget(QtGui.QWidget):
    '''
    Widget che crea l'oggetto "post", visualizzabile nella finestra principale.
    '''
    _label = None
    _textArea = None
    _testoRimosso = QtCore.pyqtSignal(int,int,int,name='testoRimosso')
    _testoAggiunto = QtCore.pyqtSignal(int,str,int,name='testoAggiunto')
    _idpost = None
    _selected = False
    def __init__(self,idpost,parent = None):
        QtGui.QWidget.__init__(self)
        self.setupUi(self)
        self._idpost = idpost
        if(self._idpost == 0):
            self._textArea.setMaximumHeight(40)
            self.setMaximumHeight(40)
            self._label.setMaximumHeight(40)
            self._label.setText("Titolo:")
        else:
            self._label.setText('id:'+str(self._idpost))
        self.horizontalLayout.addWidget(self._label)
        self.horizontalLayout.addWidget(self._textArea)
        
    def emettiRimosso(self,posizione,rimossi):
        '''
        Emette un segnale di rimozione testo.
        '''
        self.emit(QtCore.SIGNAL('testoRimosso(int,int,int)'),posizione,rimossi,self._idpost)
        
    def emettiAggiunto(self,posizione,aggiunti):
        '''
        Emette un segnale di aggiunta testo, con la posizione e i caratteri.
        '''
        self.emit(QtCore.SIGNAL('testoAggiunto(int,QString,int)'),posizione,aggiunti,self._idpost)
    
    def rimuoviTesto(self,posizione,rimossi):
        '''
        Rimuove il testo da una certa posizione e per un tot di caratteri
        '''
        print('Rimossi '+str(rimossi)+' caratteri dal post '+str(self._idpost))
        self._textArea.rimuoviTesto(posizione, rimossi)

    def aggiungiTesto(self,posizione,stringa):
        '''
        aggiunge il testo su una posizione.
        '''
        print('aggiunta la stringa "'+stringa+'" al post '+str(self._idpost))
        self._textArea.aggiungiTesto(posizione, stringa)
        
    def setupUi(self, Form):
        '''
        Imposta la finestra con gli oggetti.
        '''
        Form.setObjectName("Form")
        Form.resize(50, 40)
        self.setMaximumSize(2096,150)
        self.setMinimumSize(0, 50)
        self._label = QtGui.QLabel()
        self._textArea = PostArea.Post()
        self.horizontalLayout = QtGui.QHBoxLayout(Form)
        self.horizontalLayout.setSpacing(3)
        self.horizontalLayout.setMargin(2)
        self.horizontalLayout.setObjectName("horizontalLayout")
        
        QtCore.QObject.connect(self._textArea, QtCore.SIGNAL("testoRimosso(int,int)"),self.emettiRimosso)
        QtCore.QObject.connect(self._textArea, QtCore.SIGNAL("testoAggiunto(int,QString)"),self.emettiAggiunto)
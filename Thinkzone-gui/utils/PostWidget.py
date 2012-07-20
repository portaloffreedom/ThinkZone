'''
Created on 15/lug/2012
Widget che rappresenta un post per il programma.
I post hanno un ID che li rende unici.
@author: stengun
'''
from utils import PostArea
from PyQt4 import QtGui, QtCore

class postWidget(QtGui.QWidget):
    '''
    Widget che crea l'oggetto "post", capace di essere un post normale o una risposta a un post.
    '''
    _label = None
    _textArea = None
    _testoRimosso = QtCore.pyqtSignal(int,int,int,name='testoRimosso')
    _testoAggiunto = QtCore.pyqtSignal(int,str,int,name='testoAggiunto')
    _idpost = None
    def __init__(self,idpost,parent = None):
        QtGui.QWidget.__init__(self)
        self.setupUi(self)
        self._idpost = idpost
        if(parent != None):
            spacerItem = QtGui.QSpacerItem(40, 20, QtGui.QSizePolicy.Expanding, QtGui.QSizePolicy.Minimum)
            self.horizontalLayout.addItem(spacerItem)
        self._label.setText('id:'+str(self._idpost))
        self.horizontalLayout.addWidget(self._textArea)
        self.horizontalLayout.addWidget(self._label)
        
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
        Form.resize(363, 105)
        self.setMaximumSize(2096,150)
        self.setMinimumSize(0, 100)
        self._label = QtGui.QLabel()
        self._textArea = PostArea.Post()
        self.horizontalLayout = QtGui.QHBoxLayout(Form)
        self.horizontalLayout.setSpacing(3)
        self.horizontalLayout.setMargin(2)
        self.horizontalLayout.setObjectName("horizontalLayout")
        
        QtCore.QObject.connect(self._textArea, QtCore.SIGNAL("testoRimosso(int,int)"),self.emettiRimosso)
        QtCore.QObject.connect(self._textArea, QtCore.SIGNAL("testoAggiunto(int,QString)"),self.emettiAggiunto)
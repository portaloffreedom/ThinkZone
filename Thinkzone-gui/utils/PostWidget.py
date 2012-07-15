'''
Created on 15/lug/2012

@author: stengun
'''
from utils import postwidget, PostArea
from PyQt4 import QtGui, QtCore

try:
    _fromUtf8 = QtCore.QString.fromUtf8
except AttributeError:
    _fromUtf8 = lambda s: s

class postWidget(QtGui.QWidget):
    '''
    Widget che crea l'oggetto "post", capace di essere un post normale o una risposta a un post.
    '''
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
        self.horizontalLayout.addWidget(self._textArea)
        
    def emettiRimosso(self,posizione,rimossi):
        self.emit(QtCore.SIGNAL('testoRimosso(int,int,int)'),posizione,rimossi,self._idpost)
        
    def emettiAggiunto(self,posizione,aggiunti):
        self.emit(QtCore.SIGNAL('testoAggiunto(int,QString,int)'),posizione,aggiunti,self._idpost)
    
    def rimuoviTesto(self,posizione,rimossi):
        self._textArea.rimuoviTesto(posizione, rimossi)

    def aggiungiTesto(self,posizione,stringa):
        self._textArea.aggiungiTesto(posizione, stringa)
        
    def setupUi(self, Form):
        Form.setObjectName(_fromUtf8("Form"))
        Form.resize(363, 105)
        self.setMaximumSize(2096,150)
        self.setMinimumSize(0, 100)
        self._textArea = PostArea.Post()
        self.horizontalLayout = QtGui.QHBoxLayout(Form)
        self.horizontalLayout.setSpacing(3)
        self.horizontalLayout.setMargin(2)
        self.horizontalLayout.setObjectName(_fromUtf8("horizontalLayout"))
        
        QtCore.QObject.connect(self._textArea, QtCore.SIGNAL("testoRimosso(int,int)"),self.emettiRimosso)
        QtCore.QObject.connect(self._textArea, QtCore.SIGNAL("testoAggiunto(int,QString)"),self.emettiAggiunto)
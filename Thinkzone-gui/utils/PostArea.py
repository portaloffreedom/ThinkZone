'''
Created on 11/lug/2012

@author: stengun
'''

from PyQt4 import QtGui, QtCore

class Post(QtGui.QTextEdit):
    '''
    Classe per gestire il testo in arrivo e da inviare.
    '''
    testo = None
    
    def __init__(self,parent=None):
        QtGui.QTextEdit.__init__(self,parent)
        self._tcpSync(True)
        _testoRimosso = QtCore.pyqtSignal(int,int,name='testoRimosso')
        _testoAggiunto = QtCore.pyqtSignal(int,str,name='testoAggiunto')
    
    def _tcpSync(self,stato):
        if(stato):
            QtCore.QObject.connect(self.document(), QtCore.SIGNAL("contentsChange(int,int,int)"),self.testoCambiatoSlot)
        else:
            QtCore.QObject.disconnect(self.document(), QtCore.SIGNAL("contentsChange(int,int,int)"),self.testoCambiatoSlot)
     
    def testoCambiatoSlot(self,posizione,rimossi,aggiunti):
        if(rimossi!=0):
            #print('rimozione')
            self.emit(QtCore.SIGNAL('testoRimosso(int,int)'),posizione,rimossi)
        if(aggiunti!=0):
            self.testo = self.toPlainText()
            self.testo = self.testo[posizione:posizione+aggiunti]
            self.emit(QtCore.SIGNAL('testoAggiunto(int,QString)'),posizione,self.testo)
            posizione += aggiunti
        cursore = self.textCursor()
        cursore.setPosition(posizione)
        self.setTextCursor(cursore)
        
        
    def textUpdate(self):    
        self.setText(self.testo)
        self.update()
    
    def aggiungiTesto(self,posizione,stringa):
        print('aggiungo')
        cursore = self.textCursor()
        cursore.setPosition(posizione)
        testo1 = self.testo[:posizione]
        self.testo = testo1 + stringa + self.testo[posizione:]
        self._tcpSync(False)
        self.textUpdate()
        self._tcpSync(True)
        
    def rimuoviTesto(self,posizione,rimossi):
        print('tolgo')
        cursore = self.textCursor()
        cursore.setPosition(posizione)
        testo1 = self.testo[:posizione-rimossi]
        temp = self.testo[posizione:]
        self.testo = testo1+temp
        self._tcpSync(False)
        self.textUpdate()
        self._tcpSync(True)
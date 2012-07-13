'''
Created on 11/lug/2012

@author: stengun
'''

from PyQt4 import QtGui, QtCore

class Post(QtGui.QTextEdit):
    '''
    Classe per gestire il testo in arrivo e da inviare.
    '''
    testo = ''
    blink_cursor = 0
    
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
        self.blink_cursor = posizione
        self.setTextCursor(cursore)
        
        
    def textUpdate(self,posizione,quanti):
        self.blink_cursor = self.textCursor().position()
        if(posizione < self.blink_cursor):
            if (quanti < 0):
                distanza = (posizione - quanti - self.blink_cursor)
                if (distanza >0):
                    quanti += distanza        
            self.blink_cursor +=quanti
        cursore = self.textCursor()    
        self.setText(self.testo)
        cursore.setPosition(self.blink_cursor)
        self.update()
        self.setTextCursor(cursore)
    
    def aggiungiTesto(self,posizione,stringa):
        self.testo = self.toPlainText()
        print('aggiungo',stringa,'in posizione',posizione)
        cursore = self.textCursor()
        cursore.setPosition(posizione)
        testo1 = self.testo[:posizione]
        self.testo = testo1 + stringa + self.testo[posizione:]
        posizione = posizione + len(stringa)
        self._tcpSync(False)
        self.textUpdate(posizione,len(stringa))
        self._tcpSync(True)
        
    def rimuoviTesto(self,posizione,rimossi):
        print('tolgo ',rimossi,' caratteri dalla posizione ',posizione)
        cursore = self.textCursor()
        cursore.setPosition(posizione)
        prima = self.testo[:posizione]
        dopo = self.testo[posizione+rimossi:]
        self.testo = prima+dopo
        self._tcpSync(False)
        self.textUpdate(posizione,-rimossi)
        self._tcpSync(True)
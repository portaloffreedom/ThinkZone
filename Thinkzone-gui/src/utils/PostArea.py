'''
Mini widget per gestire la modifica concorrenziale di una textarea.
Eredita da *QTextEdit* e sviluppa in modo semplice i metodi per consentire la
scrittura a più utenti senza modificare puntatori all'utente locale.
'''

from PyQt4 import QtGui, QtCore

class Post(QtGui.QTextEdit):
    '''
    Classe che rappresenta la textArea del post
    '''
    testo = ''
    blink_cursor = 0
    
    def __init__(self,parent=None):
        QtGui.QTextEdit.__init__(self,parent)
        self.setMinimumHeight(40)
        self._tcpSync(True)
        _testoRimosso = QtCore.pyqtSignal(int,int,name='testoRimosso')
        _testoAggiunto = QtCore.pyqtSignal(int,str,name='testoAggiunto')
    
    def _tcpSync(self,stato):
        if(stato):
            QtCore.QObject.connect(self.document(), QtCore.SIGNAL("contentsChange(int,int,int)"),self.testoCambiatoSlot)
        else:
            QtCore.QObject.disconnect(self.document(), QtCore.SIGNAL("contentsChange(int,int,int)"),self.testoCambiatoSlot)
     
    def testoCambiatoSlot(self,posizione,rimossi,aggiunti):
        '''
        Metodo slot, viene chiamato quando il testo nella text area cambia.
        Calcola le varie posizioni e capisce cosa è stato aggiunto o rimosso. Emette i segnali in base a 
        cosa è successo.
        '''
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
        '''
        Metodo chiamato quando la textArea è da aggiornare.
        Fa in modo che il cursore dell'utente locale lampeggi sempre nello stesso punto, 
        senza spostarsi.
        '''
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
        '''
        Aggiunge una stringa di testo nella textArea.
        '''
        self.testo = self.toPlainText()
        #print('aggiungo',stringa,'in posizione',posizione)
        cursore = self.textCursor()
        cursore.setPosition(posizione)
        testo1 = self.testo[:posizione]
        self.testo = testo1 + stringa + self.testo[posizione:]
        posizione = posizione + len(stringa)
        self._tcpSync(False)
        self.textUpdate(posizione,len(stringa))
        self._tcpSync(True)
        
    def rimuoviTesto(self,posizione,rimossi):
        '''
        Rimuove un numero arbitrario di caratteri dalla textarea.
        '''
        #print('tolgo ',rimossi,' caratteri dalla posizione ',posizione)
        self.testo = self.toPlainText()
        cursore = self.textCursor()
        cursore.setPosition(posizione)
        prima = self.testo[:posizione]
        dopo = self.testo[posizione+rimossi:]
        self.testo = prima+dopo
        self._tcpSync(False)
        self.textUpdate(posizione,-rimossi)
        self._tcpSync(True)
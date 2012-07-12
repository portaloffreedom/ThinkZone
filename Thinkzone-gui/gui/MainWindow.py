'''
Created on 11/lug/2012
Finestra principale che si occupa di connettere le componenti della GUI e i widget personalizzati.
@author: stengun
'''
import sys
from utils import PostArea
from rete import Comunicazione
from gui import mwind
from PyQt4 import QtGui, QtCore

class FinestraPrincipale(QtGui.QMainWindow, mwind.Ui_finestraprincipale):
    '''
    Classe per la finestra principale.
    '''
    _textextended = None
    _connettore = None
    def __init__(self, parent = None):
        QtGui.QMainWindow.__init__(self, parent)
        self._textextended = PostArea.Post()
        self.ui = mwind.Ui_finestraprincipale()
        self.setupUi(self)
        self.scroll_layout_2.addWidget(self._textextended)
        self._connettore = Comunicazione.comunicatore()
        QtCore.QObject.connect(self.bottone_connetti,QtCore.SIGNAL('pressed()'), self.connetti)
        QtCore.QObject.connect(self._textextended,QtCore.SIGNAL('testoRimosso(int,int)'), self.dati_rimossi)
        QtCore.QObject.connect(self._textextended,QtCore.SIGNAL('testoAggiunto(int,QString)'), self.dati_aggiunti)
    
    def dati_rimossi(self,posizione,rimossi):
        #print('ci sono')
        self._connettore.spedisci_rimozione(posizione,rimossi)
    
    def dati_aggiunti(self,posizione,aggiunti):
        self._connettore.spedisci_aggiunta(posizione,aggiunti)
    
    def connetti(self):
        porta = self.testo_porta.text()
        porta = porta.encode()
        porta = int(porta)
        self._connettore.connetti(self.testo_host.text(), porta,self.testo_nickname.text())
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('rimozione(int,int)'),self._textextended.rimuoviTesto)
        QtCore.QObject.connect(self._connettore,QtCore.SIGNAL('aggiunta(int,QString)'),self._textextended.aggiungiTesto)
        self._connettore.start()

        
if __name__ == '__main__':
    app = QtGui.QApplication(sys.argv)
    finestra = FinestraPrincipale()
    finestra.show()
    app.exec()
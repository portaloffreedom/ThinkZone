'''
Created on 12/lug/2012
Finestra di Dialog per la schermata di ABOUT.
Aggiunge poco e nulla al costruttore generato da QT Designer.

@author: stengun
'''
from gui import about
from PyQt4 import QtGui, QtCore
class aboutDial(QtGui.QDialog, about.Ui_Dialog):
    '''
    Costruisce una finestra di dialogo ABOUT per Thinkzone.
    '''


    def __init__(self,parent = None):
        QtGui.QDialog.__init__(self,parent)
        self.ui = about.Ui_Dialog()
        self.setupUi(self)
        QtCore.QObject.connect(self.buttonChiudi, QtCore.SIGNAL('pressed()'),self.close)

        
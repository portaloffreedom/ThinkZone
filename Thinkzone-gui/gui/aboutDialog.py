'''
Created on 12/lug/2012

@author: stengun
'''
from gui import about
from PyQt4 import QtGui, QtCore
class aboutDial(QtGui.QDialog, about.Ui_Dialog):
    '''
    classdocs
    '''


    def __init__(self,parent = None):
        QtGui.QDialog.__init__(self)
        about.Ui_Dialog.setupUi(self,self)
        QtCore.QObject.connect(self.buttonChiudi, QtCore.SIGNAL('pressed()'),self.close)
        '''
        Constructor
        '''
        
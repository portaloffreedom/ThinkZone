'''
Unittest per verificare le impostazioni delle finestre in presenza di input non ammessi.
@author: stengun
'''
import unittest,sys
from PyQt4.QtGui import QApplication
from PyQt4.QtCore import Qt
from PyQt4.QtTest import QTest
from rete import Comunicazione
from gui import aboutDialog,loginDialog,MainWindow

class dummyServer():
    pass

class testWidgets(unittest.TestCase):
    
    def setUp(self):
        self.app = QApplication(sys.argv)
        self._principale = MainWindow.mainwindow("0.0.0")
        self._about = aboutDialog.aboutDial(self._principale)
        self._loginwidget = loginDialog.Login(self._principale)
    
    def testLoginWidget(self):
        self._loginwidget.usernameEdit.setText("testing")
        self._loginwidget.passwordEdit.setText("testing")
        self._loginwidget.hostEdit.setText("192.168.0.42")
        
        self._loginwidget.portaEdit.setText("4242")
        self.assertGreater(int(self._loginwidget.portaEdit.text()), 0, "Testing valore porta")
        self._loginwidget.portaEdit.setText("-4242")
        self.assertGreater(int(self._loginwidget.portaEdit.text()), 0, "Testing valore porta")
        self._loginwidget.portaEdit.setText("a32o")
        self.assertGreater(int(self._loginwidget.portaEdit.text()), 0, "Testing valore porta")
        
        registrabutton = self._loginwidget.buttonRegister
        connettibutton = self._loginwidget.buttonConnect
        #sequenza sballata
        QTest.mouseClick(connettibutton,Qt.LeftButton)
        QTest.mouseClick(registrabutton,Qt.LeftButton)
        #sequenza corretta
        QTest.mouseClick(registrabutton,Qt.LeftButton)
        QTest.mouseClick(connettibutton,Qt.LeftButton)

        
    def testAboutDialog(self):
        pass
    
    def testPostAdd(self):
        pass


def suite():
    suite = unittest.TestSuite()
    suite.addTest(testWidgets)
    return suite


if __name__ == "__main__":
    #import sys;sys.argv = ['', 'Test.testLogin']
    unittest.main()
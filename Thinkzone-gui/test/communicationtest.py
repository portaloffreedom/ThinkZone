'''
UNITTEST per testare i metodi di registrazione, connessione e disconnessione del server.
@author: stengun
'''
import unittest
from rete import Comunicazione
from gui import MainWindow

class TestLogin(unittest.TestCase):
    
    login = None
    hostname = "192.168.0.42"
    porta = 22
    nome = "testingoo"
    password = "testing"
    
    def setUp(self):
        self.login = Comunicazione.comunicatore()
        self.login._postPlexer = MainWindow.PostPlexer(self.login)
    
#    def testRegister(self):
#        self.assertFalse(self.login.registrati(self.hostname, 'abc', self.nome, self.password),"Testo la registrazione")
#        self.assertTrue(self.login.registrati(self.hostname, self.porta, self.nome, self.password),"Testo la registrazione")
        
    def testLogin(self):     
        self.assertTrue(self.login.connetti(self.hostname, self.porta, self.nome, self.password),"Testando connessione")
        if(self.login._postPlexer._userID != None):
            self.assertGreater(self.login._postPlexer._userID, 0, "testando l'user ID")
    
    def testDisconnessione(self):
        self.login.disconnetti()
        if(self.login._receive_thread != None):
            self.assertFalse(self.login._receive_thread.isRunning(),"Testando l'uccisione del thread")
        self.assertFalse(self.login.isRunning(),"Verifico se il comunicatore sta ancora girando")
         

if __name__ == "__main__":
    #import sys;sys.argv = ['', 'Test.testName']
    unittest.main()
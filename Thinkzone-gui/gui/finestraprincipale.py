# -*- coding: utf-8 -*-

# Form implementation generated from reading ui file 'mainwindow.ui'
#
# Created: Thu Jul 12 09:49:05 2012
#      by: PyQt4 UI code generator 4.9.4
#
# WARNING! All changes made in this file will be lost!

from PyQt4 import QtCore, QtGui

try:
    _fromUtf8 = QtCore.QString.fromUtf8
except AttributeError:
    _fromUtf8 = lambda s: s

class Ui_finestraprincipale(object):
    def setupUi(self, finestraprincipale):
        finestraprincipale.setObjectName(_fromUtf8("finestraprincipale"))
        finestraprincipale.setEnabled(True)
        finestraprincipale.resize(758, 697)
        self.wid_centr = QtGui.QWidget(finestraprincipale)
        self.wid_centr.setObjectName(_fromUtf8("wid_centr"))
        self.verticalLayout_4 = QtGui.QVBoxLayout(self.wid_centr)
        self.verticalLayout_4.setObjectName(_fromUtf8("verticalLayout_4"))
        self.frame_3 = QtGui.QFrame(self.wid_centr)
        self.frame_3.setFrameShape(QtGui.QFrame.StyledPanel)
        self.frame_3.setFrameShadow(QtGui.QFrame.Raised)
        self.frame_3.setObjectName(_fromUtf8("frame_3"))
        self.verticalLayout_3 = QtGui.QVBoxLayout(self.frame_3)
        self.verticalLayout_3.setObjectName(_fromUtf8("verticalLayout_3"))
        self.frame_2 = QtGui.QFrame(self.frame_3)
        self.frame_2.setEnabled(True)
        self.frame_2.setFrameShape(QtGui.QFrame.StyledPanel)
        self.frame_2.setFrameShadow(QtGui.QFrame.Raised)
        self.frame_2.setObjectName(_fromUtf8("frame_2"))
        self.horizontalLayout_2 = QtGui.QHBoxLayout(self.frame_2)
        self.horizontalLayout_2.setObjectName(_fromUtf8("horizontalLayout_2"))
        self.frame_indirizzo = QtGui.QFrame(self.frame_2)
        self.frame_indirizzo.setObjectName(_fromUtf8("frame_indirizzo"))
        self.verticalLayout = QtGui.QVBoxLayout(self.frame_indirizzo)
        self.verticalLayout.setObjectName(_fromUtf8("verticalLayout"))
        self.label_indirizzo = QtGui.QLabel(self.frame_indirizzo)
        self.label_indirizzo.setObjectName(_fromUtf8("label_indirizzo"))
        self.verticalLayout.addWidget(self.label_indirizzo)
        self.testo_host = QtGui.QLineEdit(self.frame_indirizzo)
        self.testo_host.setEnabled(True)
        self.testo_host.setText(_fromUtf8(""))
        self.testo_host.setMaxLength(255)
        self.testo_host.setObjectName(_fromUtf8("testo_host"))
        self.verticalLayout.addWidget(self.testo_host)
        self.horizontalLayout_2.addWidget(self.frame_indirizzo)
        self.frame_porta = QtGui.QFrame(self.frame_2)
        self.frame_porta.setMaximumSize(QtCore.QSize(100, 16777215))
        self.frame_porta.setObjectName(_fromUtf8("frame_porta"))
        self.verticalLayout_2 = QtGui.QVBoxLayout(self.frame_porta)
        self.verticalLayout_2.setObjectName(_fromUtf8("verticalLayout_2"))
        self.label_porta = QtGui.QLabel(self.frame_porta)
        self.label_porta.setObjectName(_fromUtf8("label_porta"))
        self.verticalLayout_2.addWidget(self.label_porta)
        self.testo_porta = QtGui.QLineEdit(self.frame_porta)
        self.testo_porta.setInputMethodHints(QtCore.Qt.ImhDigitsOnly|QtCore.Qt.ImhFormattedNumbersOnly)
        self.testo_porta.setMaxLength(5)
        self.testo_porta.setObjectName(_fromUtf8("testo_porta"))
        self.verticalLayout_2.addWidget(self.testo_porta)
        self.horizontalLayout_2.addWidget(self.frame_porta)
        self.widget_conn = QtGui.QWidget(self.frame_2)
        self.widget_conn.setLayoutDirection(QtCore.Qt.RightToLeft)
        self.widget_conn.setObjectName(_fromUtf8("widget_conn"))
        self.verticalLayout_5 = QtGui.QVBoxLayout(self.widget_conn)
        self.verticalLayout_5.setMargin(0)
        self.verticalLayout_5.setObjectName(_fromUtf8("verticalLayout_5"))
        self.testo_nickname = QtGui.QLineEdit(self.widget_conn)
        self.testo_nickname.setObjectName(_fromUtf8("testo_nickname"))
        self.verticalLayout_5.addWidget(self.testo_nickname)
        self.bottone_connetti = QtGui.QPushButton(self.widget_conn)
        self.bottone_connetti.setObjectName(_fromUtf8("bottone_connetti"))
        self.verticalLayout_5.addWidget(self.bottone_connetti)
        self.horizontalLayout_2.addWidget(self.widget_conn)
        self.verticalLayout_3.addWidget(self.frame_2)
        self.frame_coll = QtGui.QFrame(self.frame_3)
        self.frame_coll.setFrameShape(QtGui.QFrame.StyledPanel)
        self.frame_coll.setFrameShadow(QtGui.QFrame.Raised)
        self.frame_coll.setObjectName(_fromUtf8("frame_coll"))
        self.verticalLayout_6 = QtGui.QVBoxLayout(self.frame_coll)
        self.verticalLayout_6.setObjectName(_fromUtf8("verticalLayout_6"))
        self.scrollArea = QtGui.QScrollArea(self.frame_coll)
        self.scrollArea.setWidgetResizable(True)
        self.scrollArea.setObjectName(_fromUtf8("scrollArea"))
        self.scroll_widget = QtGui.QWidget()
        self.scroll_widget.setGeometry(QtCore.QRect(0, 0, 698, 453))
        self.scroll_widget.setObjectName(_fromUtf8("scroll_widget"))
        self.scroll_layout_2 = QtGui.QGridLayout(self.scroll_widget)
        self.scroll_layout_2.setMargin(0)
        self.scroll_layout_2.setSpacing(0)
        self.scroll_layout_2.setObjectName(_fromUtf8("scroll_layout_2"))
        self.scrollArea.setWidget(self.scroll_widget)
        self.verticalLayout_6.addWidget(self.scrollArea)
        self.verticalLayout_3.addWidget(self.frame_coll)
        self.verticalLayout_4.addWidget(self.frame_3)
        finestraprincipale.setCentralWidget(self.wid_centr)
        self.menubar = QtGui.QMenuBar(finestraprincipale)
        self.menubar.setGeometry(QtCore.QRect(0, 0, 758, 30))
        self.menubar.setObjectName(_fromUtf8("menubar"))
        self.menuFail = QtGui.QMenu(self.menubar)
        self.menuFail.setObjectName(_fromUtf8("menuFail"))
        finestraprincipale.setMenuBar(self.menubar)
        self.statusbar = QtGui.QStatusBar(finestraprincipale)
        self.statusbar.setObjectName(_fromUtf8("statusbar"))
        finestraprincipale.setStatusBar(self.statusbar)
        self.toolBar = QtGui.QToolBar(finestraprincipale)
        self.toolBar.setObjectName(_fromUtf8("toolBar"))
        finestraprincipale.addToolBar(QtCore.Qt.TopToolBarArea, self.toolBar)
        self.actionEsci = QtGui.QAction(finestraprincipale)
        self.actionEsci.setObjectName(_fromUtf8("actionEsci"))
        self.actionConnetti = QtGui.QAction(finestraprincipale)
        self.actionConnetti.setObjectName(_fromUtf8("actionConnetti"))
        self.menuFail.addAction(self.actionEsci)
        self.menubar.addAction(self.menuFail.menuAction())
        self.label_indirizzo.setBuddy(self.testo_host)
        self.label_porta.setBuddy(self.testo_porta)

        self.retranslateUi(finestraprincipale)
        QtCore.QObject.connect(self.actionEsci, QtCore.SIGNAL(_fromUtf8("triggered()")), finestraprincipale.close)
        QtCore.QMetaObject.connectSlotsByName(finestraprincipale)

    def retranslateUi(self, finestraprincipale):
        finestraprincipale.setWindowTitle(QtGui.QApplication.translate("finestraprincipale", "Netcaz", None, QtGui.QApplication.UnicodeUTF8))
        self.label_indirizzo.setText(QtGui.QApplication.translate("finestraprincipale", "indirizzo", None, QtGui.QApplication.UnicodeUTF8))
        self.testo_host.setPlaceholderText(QtGui.QApplication.translate("finestraprincipale", "hostname", None, QtGui.QApplication.UnicodeUTF8))
        self.label_porta.setText(QtGui.QApplication.translate("finestraprincipale", "porta", None, QtGui.QApplication.UnicodeUTF8))
        self.testo_porta.setPlaceholderText(QtGui.QApplication.translate("finestraprincipale", "4000", None, QtGui.QApplication.UnicodeUTF8))
        self.testo_nickname.setPlaceholderText(QtGui.QApplication.translate("finestraprincipale", "Nickname", None, QtGui.QApplication.UnicodeUTF8))
        self.bottone_connetti.setText(QtGui.QApplication.translate("finestraprincipale", "Connetti", None, QtGui.QApplication.UnicodeUTF8))
        self.menuFail.setTitle(QtGui.QApplication.translate("finestraprincipale", "fail", None, QtGui.QApplication.UnicodeUTF8))
        self.toolBar.setWindowTitle(QtGui.QApplication.translate("finestraprincipale", "toolBar", None, QtGui.QApplication.UnicodeUTF8))
        self.actionEsci.setText(QtGui.QApplication.translate("finestraprincipale", "esci", None, QtGui.QApplication.UnicodeUTF8))
        self.actionConnetti.setText(QtGui.QApplication.translate("finestraprincipale", "Connetti", None, QtGui.QApplication.UnicodeUTF8))
        self.actionConnetti.setToolTip(QtGui.QApplication.translate("finestraprincipale", "connettiti", None, QtGui.QApplication.UnicodeUTF8))


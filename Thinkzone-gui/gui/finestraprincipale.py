# -*- coding: utf-8 -*-

# Form implementation generated from reading ui file 'conversations.ui'
#
# Created: Sun Jul 15 22:12:50 2012
#      by: PyQt4 UI code generator 4.9.4
#
# WARNING! All changes made in this file will be lost!

from PyQt4 import QtCore, QtGui

try:
    _fromUtf8 = QtCore.QString.fromUtf8
except AttributeError:
    _fromUtf8 = lambda s: s

class Ui_MainWindow(object):
    def setupUi(self, MainWindow):
        MainWindow.setObjectName(_fromUtf8("MainWindow"))
        MainWindow.resize(1117, 869)
        icon = QtGui.QIcon()
        icon.addPixmap(QtGui.QPixmap(_fromUtf8(":/images/gelatino_icon2.png")), QtGui.QIcon.Normal, QtGui.QIcon.Off)
        MainWindow.setWindowIcon(icon)
        self.centralwidget = QtGui.QWidget(MainWindow)
        self.centralwidget.setObjectName(_fromUtf8("centralwidget"))
        self.horizontalLayout = QtGui.QHBoxLayout(self.centralwidget)
        self.horizontalLayout.setSpacing(2)
        self.horizontalLayout.setMargin(2)
        self.horizontalLayout.setObjectName(_fromUtf8("horizontalLayout"))
        self.scrollArea = QtGui.QScrollArea(self.centralwidget)
        self.scrollArea.setWidgetResizable(True)
        self.scrollArea.setObjectName(_fromUtf8("scrollArea"))
        self.scrollAreaWidgetContents = QtGui.QWidget()
        self.scrollAreaWidgetContents.setGeometry(QtCore.QRect(0, 0, 849, 725))
        self.scrollAreaWidgetContents.setObjectName(_fromUtf8("scrollAreaWidgetContents"))
        self.layoutTextarea = QtGui.QVBoxLayout(self.scrollAreaWidgetContents)
        self.layoutTextarea.setObjectName(_fromUtf8("layoutTextarea"))
        self.scrollArea.setWidget(self.scrollAreaWidgetContents)
        self.horizontalLayout.addWidget(self.scrollArea)
        MainWindow.setCentralWidget(self.centralwidget)
        self.menubar = QtGui.QMenuBar(MainWindow)
        self.menubar.setGeometry(QtCore.QRect(0, 0, 1117, 30))
        self.menubar.setObjectName(_fromUtf8("menubar"))
        self.menuFile = QtGui.QMenu(self.menubar)
        self.menuFile.setObjectName(_fromUtf8("menuFile"))
        self.menuAiuto = QtGui.QMenu(self.menubar)
        self.menuAiuto.setObjectName(_fromUtf8("menuAiuto"))
        self.menuVisualizza = QtGui.QMenu(self.menubar)
        self.menuVisualizza.setObjectName(_fromUtf8("menuVisualizza"))
        MainWindow.setMenuBar(self.menubar)
        self.statusbar = QtGui.QStatusBar(MainWindow)
        self.statusbar.setObjectName(_fromUtf8("statusbar"))
        MainWindow.setStatusBar(self.statusbar)
        self.dockLista = QtGui.QDockWidget(MainWindow)
        sizePolicy = QtGui.QSizePolicy(QtGui.QSizePolicy.Preferred, QtGui.QSizePolicy.Preferred)
        sizePolicy.setHorizontalStretch(0)
        sizePolicy.setVerticalStretch(0)
        sizePolicy.setHeightForWidth(self.dockLista.sizePolicy().hasHeightForWidth())
        self.dockLista.setSizePolicy(sizePolicy)
        self.dockLista.setAllowedAreas(QtCore.Qt.LeftDockWidgetArea|QtCore.Qt.RightDockWidgetArea)
        self.dockLista.setObjectName(_fromUtf8("dockLista"))
        self.dockWidgetContents = QtGui.QWidget()
        self.dockWidgetContents.setObjectName(_fromUtf8("dockWidgetContents"))
        self.verticalLayout_2 = QtGui.QVBoxLayout(self.dockWidgetContents)
        self.verticalLayout_2.setSpacing(3)
        self.verticalLayout_2.setMargin(0)
        self.verticalLayout_2.setObjectName(_fromUtf8("verticalLayout_2"))
        self.listaConvers = QtGui.QListView(self.dockWidgetContents)
        self.listaConvers.setFrameShape(QtGui.QFrame.StyledPanel)
        self.listaConvers.setFrameShadow(QtGui.QFrame.Plain)
        self.listaConvers.setLineWidth(1)
        self.listaConvers.setMidLineWidth(0)
        self.listaConvers.setHorizontalScrollBarPolicy(QtCore.Qt.ScrollBarAlwaysOff)
        self.listaConvers.setEditTriggers(QtGui.QAbstractItemView.DoubleClicked|QtGui.QAbstractItemView.EditKeyPressed|QtGui.QAbstractItemView.SelectedClicked)
        self.listaConvers.setAlternatingRowColors(True)
        self.listaConvers.setSelectionBehavior(QtGui.QAbstractItemView.SelectRows)
        self.listaConvers.setResizeMode(QtGui.QListView.Adjust)
        self.listaConvers.setObjectName(_fromUtf8("listaConvers"))
        self.verticalLayout_2.addWidget(self.listaConvers)
        self.widgetPartecipanti = QtGui.QWidget(self.dockWidgetContents)
        self.widgetPartecipanti.setObjectName(_fromUtf8("widgetPartecipanti"))
        self.layoutPartecipanti = QtGui.QGridLayout(self.widgetPartecipanti)
        self.layoutPartecipanti.setMargin(0)
        self.layoutPartecipanti.setObjectName(_fromUtf8("layoutPartecipanti"))
        self.verticalLayout_2.addWidget(self.widgetPartecipanti)
        self.dockLista.setWidget(self.dockWidgetContents)
        MainWindow.addDockWidget(QtCore.Qt.DockWidgetArea(1), self.dockLista)
        self.dock_strumconv = QtGui.QDockWidget(MainWindow)
        self.dock_strumconv.setObjectName(_fromUtf8("dock_strumconv"))
        self.dockWidgetContents_2 = QtGui.QWidget()
        self.dockWidgetContents_2.setObjectName(_fromUtf8("dockWidgetContents_2"))
        self.horizontalLayout_2 = QtGui.QHBoxLayout(self.dockWidgetContents_2)
        self.horizontalLayout_2.setObjectName(_fromUtf8("horizontalLayout_2"))
        spacerItem = QtGui.QSpacerItem(40, 20, QtGui.QSizePolicy.Expanding, QtGui.QSizePolicy.Minimum)
        self.horizontalLayout_2.addItem(spacerItem)
        self.titoloEdit = QtGui.QLineEdit(self.dockWidgetContents_2)
        self.titoloEdit.setObjectName(_fromUtf8("titoloEdit"))
        self.horizontalLayout_2.addWidget(self.titoloEdit)
        self.buttonCrea = QtGui.QPushButton(self.dockWidgetContents_2)
        self.buttonCrea.setObjectName(_fromUtf8("buttonCrea"))
        self.horizontalLayout_2.addWidget(self.buttonCrea)
        self.buttonElimina = QtGui.QPushButton(self.dockWidgetContents_2)
        self.buttonElimina.setEnabled(False)
        self.buttonElimina.setObjectName(_fromUtf8("buttonElimina"))
        self.horizontalLayout_2.addWidget(self.buttonElimina)
        self.dock_strumconv.setWidget(self.dockWidgetContents_2)
        MainWindow.addDockWidget(QtCore.Qt.DockWidgetArea(4), self.dock_strumconv)
        self.actionLogin = QtGui.QAction(MainWindow)
        self.actionLogin.setObjectName(_fromUtf8("actionLogin"))
        self.actionEsci = QtGui.QAction(MainWindow)
        self.actionEsci.setObjectName(_fromUtf8("actionEsci"))
        self.actionInformazioni_su = QtGui.QAction(MainWindow)
        self.actionInformazioni_su.setObjectName(_fromUtf8("actionInformazioni_su"))
        self.actionSostienici = QtGui.QAction(MainWindow)
        self.actionSostienici.setObjectName(_fromUtf8("actionSostienici"))
        self.actionLista_conversazioni = QtGui.QAction(MainWindow)
        self.actionLista_conversazioni.setCheckable(True)
        self.actionLista_conversazioni.setChecked(True)
        self.actionLista_conversazioni.setObjectName(_fromUtf8("actionLista_conversazioni"))
        self.actionStrumenti_conversazione = QtGui.QAction(MainWindow)
        self.actionStrumenti_conversazione.setCheckable(True)
        self.actionStrumenti_conversazione.setChecked(True)
        self.actionStrumenti_conversazione.setObjectName(_fromUtf8("actionStrumenti_conversazione"))
        self.menuFile.addAction(self.actionLogin)
        self.menuFile.addAction(self.actionEsci)
        self.menuAiuto.addAction(self.actionInformazioni_su)
        self.menuAiuto.addAction(self.actionSostienici)
        self.menuVisualizza.addAction(self.actionLista_conversazioni)
        self.menuVisualizza.addAction(self.actionStrumenti_conversazione)
        self.menubar.addAction(self.menuFile.menuAction())
        self.menubar.addAction(self.menuVisualizza.menuAction())
        self.menubar.addAction(self.menuAiuto.menuAction())

        self.retranslateUi(MainWindow)
        QtCore.QObject.connect(self.actionEsci, QtCore.SIGNAL(_fromUtf8("triggered()")), MainWindow.close)
        QtCore.QObject.connect(self.actionLista_conversazioni, QtCore.SIGNAL(_fromUtf8("toggled(bool)")), self.dockLista.setVisible)
        QtCore.QObject.connect(self.actionStrumenti_conversazione, QtCore.SIGNAL(_fromUtf8("toggled(bool)")), self.dock_strumconv.setVisible)
        QtCore.QObject.connect(self.dockLista, QtCore.SIGNAL(_fromUtf8("visibilityChanged(bool)")), self.actionLista_conversazioni.setChecked)
        QtCore.QObject.connect(self.dock_strumconv, QtCore.SIGNAL(_fromUtf8("visibilityChanged(bool)")), self.actionStrumenti_conversazione.setChecked)
        QtCore.QMetaObject.connectSlotsByName(MainWindow)

    def retranslateUi(self, MainWindow):
        MainWindow.setWindowTitle(QtGui.QApplication.translate("MainWindow", "MainWindow", None, QtGui.QApplication.UnicodeUTF8))
        self.menuFile.setTitle(QtGui.QApplication.translate("MainWindow", "File", None, QtGui.QApplication.UnicodeUTF8))
        self.menuAiuto.setTitle(QtGui.QApplication.translate("MainWindow", "Aiuto!", None, QtGui.QApplication.UnicodeUTF8))
        self.menuVisualizza.setTitle(QtGui.QApplication.translate("MainWindow", "Visualizza", None, QtGui.QApplication.UnicodeUTF8))
        self.titoloEdit.setPlaceholderText(QtGui.QApplication.translate("MainWindow", "titolo del nuovo post", None, QtGui.QApplication.UnicodeUTF8))
        self.buttonCrea.setText(QtGui.QApplication.translate("MainWindow", "Crea", None, QtGui.QApplication.UnicodeUTF8))
        self.buttonElimina.setText(QtGui.QApplication.translate("MainWindow", "Elimina", None, QtGui.QApplication.UnicodeUTF8))
        self.actionLogin.setText(QtGui.QApplication.translate("MainWindow", "Login", None, QtGui.QApplication.UnicodeUTF8))
        self.actionEsci.setText(QtGui.QApplication.translate("MainWindow", "Esci", None, QtGui.QApplication.UnicodeUTF8))
        self.actionInformazioni_su.setText(QtGui.QApplication.translate("MainWindow", "Informazioni su Thinkzone", None, QtGui.QApplication.UnicodeUTF8))
        self.actionSostienici.setText(QtGui.QApplication.translate("MainWindow", "Sostienici", None, QtGui.QApplication.UnicodeUTF8))
        self.actionLista_conversazioni.setText(QtGui.QApplication.translate("MainWindow", "Lista conversazioni", None, QtGui.QApplication.UnicodeUTF8))
        self.actionStrumenti_conversazione.setText(QtGui.QApplication.translate("MainWindow", "Strumenti conversazione", None, QtGui.QApplication.UnicodeUTF8))

import immagini_rc

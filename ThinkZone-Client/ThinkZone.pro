#-------------------------------------------------
#
# Project created by QtCreator 2012-05-24T20:35:26
#
#-------------------------------------------------

QT += core gui
QT += network thread

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

TARGET = ThinkZone
TEMPLATE = app


SOURCES += main.cpp\
        mainwindow.cpp \
    inputtesto.cpp \
    superstring.cpp \
    comunicator.cpp \
    database.cpp \
    login.cpp

HEADERS  += mainwindow.h \
    inputtesto.h \
    superstring.h \
    comunicator.h \
    database.h \
    login.h

FORMS    +=

#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QtGui>
#include <QtNetwork>
#include "inputtesto.h"
#include "comunicator.h"

class MainWindow : public QWidget
{
    Q_OBJECT
    
public:
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();
    
private:
    SuperString canale;
    inputTesto input;
    inputTesto output;
//    QTcpServer *server;
//    QTcpSocket *connessione;
    Comunicator *comunicator;

//    int readInt();

signals:
//    void testoAggiunto(int position, QString *stringa);
//    void testoRimosso(int position, int howmany);

//private slots:
//    void newConnection();
//    void serve();

public slots:
    void StartMainWindow(QString *username);


};

#endif // MAINWINDOW_H

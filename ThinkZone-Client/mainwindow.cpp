#include "mainwindow.h"
#include <iostream>
using namespace std;

MainWindow::MainWindow(QWidget *parent) :
//    QMainWindow(parent)
    QWidget(parent)
{
//    comunicator = new Comunicator(0,QHostAddress("192.168.0.42"));
    comunicator = new Comunicator();
//    comunicator->connect_to_host();

    input.setTextSource(&canale);
    output.setTextSource(&canale);
    output.setReadOnly(true);

//    server = new QTcpServer(this);
//    //connect(server,SIGNAL(newConnection()),this,SLOT(newConnection()));
//    if (!server->listen(QHostAddress::LocalHost,4242)){
//        //TODO
//        cerr<<"errore server non partito!"<<endl;
//    }
//    server->waitForNewConnection(-1);
//    connessione = server->nextPendingConnection();
//    newConnection();
//    connect(this,SIGNAL(testoRimosso(int,int)),&output,SLOT(rimuoviTesto(int,int)));
//    connect(this,SIGNAL(testoAggiunto(int,QString*)),&output,SLOT(aggiungiTesto(int,QString*)));


//    comunicator->connect(comunicator,SIGNAL(testoRimosso(int,int)),&output,SLOT(rimuoviTesto(int,int)));
//    comunicator->connect(comunicator,SIGNAL(testoAggiunto(int,QString*)),&output,SLOT(aggiungiTesto(int,QString*)));
//    comunicator->connect(comunicator,SIGNAL(testoAggiunto(int,QChar)),&output,SLOT(aggiungiTesto(int,QChar)));


//    connect(input.document(),SIGNAL(contentsChange(int,int,int)),&input,SLOT(testoCambiato_Slot(int,int,int)));
    input.setSincTCP(true);
//    connect(&input,SIGNAL(testoCambiato()) , &output, SLOT(textUpdate()   ) );
    connect(&input,SIGNAL(testoAggiunto(int,QString*)),comunicator,SLOT(aggiungiTesto(int,QString*)));
    connect(&input,SIGNAL(testoRimosso(int,int)),comunicator,SLOT(rimuoviTesto(int,int)));


    comunicator->connect(comunicator,SIGNAL(testoRimosso(int,int)),&input,SLOT(rimuoviTesto(int,int)));
    comunicator->connect(comunicator,SIGNAL(testoAggiunto(int,QString*)),&input,SLOT(aggiungiTesto(int,QString*)));
    comunicator->connect(comunicator,SIGNAL(testoAggiunto(int,QChar)),&input,SLOT(aggiungiTesto(int,QChar)));


    QGridLayout *layout = new QGridLayout;
    layout->addWidget(&input, 0, 0, Qt::AlignCenter);
    layout->addWidget(&output, 1, 0, Qt::AlignCenter);

    this->setLayout(layout);
    this->setWindowTitle("ThinkZone");
}

MainWindow::~MainWindow()
{

}

void MainWindow::StartMainWindow(QString *username)
{


    comunicator->run(username);

    QString windowTitle("ThinkZone - ");
    windowTitle.append(comunicator->getUserName());
    this->setWindowTitle(windowTitle);
    this->show();
}


/*
void MainWindow::newConnection()
{
    server->close();
//    connessione = server->nextPendingConnection();
//    server->incomingConnection();
    connect(connessione, SIGNAL(disconnected()),
            connessione, SLOT(deleteLater()));
    connect(connessione, SIGNAL(readyRead()),
            this,SLOT(serve()));
}

int MainWindow::readInt()
{
    int position_I;
    const int buffSize = 64;
    char buf[buffSize];

    connessione->readLine(buf,buffSize);
    QString position_S(buf);
    position_I = position_S.toInt();

    return position_I;

}

void MainWindow::serve()
{
    //TODO manda al visualizzatore u.u
    int posizione;
    char todo;
    connessione->getChar(&todo);
    posizione = readInt();
    switch (todo) {
    case 'A': {
        const int buffSize = 64;
        char buf[buffSize];

        connessione->waitForReadyRead();
        connessione->readLine(buf,buffSize);
        QString *addString = new QString(buf);
        addString->chop(1);

        testoAggiunto(posizione,addString);
        break;
    }
    case 'D': {
        connessione->waitForReadyRead();
        int howmanyrem = readInt();
        testoRimosso(posizione,howmanyrem);
        break;
    }
    default:
        cerr<<"che cazzo sta succedendo?"<<endl;
        break;
    }
}
*/

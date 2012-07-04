#ifndef COMUNICATOR_H
#define COMUNICATOR_H

#include <QObject>
#include <QTcpSocket>
#include <QHostAddress>
#include <QThread>
#include "database.h"

#define STD_PORT 4242
#define STD_ADDRESS QHostAddress::LocalHost

class Comunicator : public QThread
{
    Q_OBJECT
public:
    explicit Comunicator(QObject *parent = 0, const QHostAddress &address = STD_ADDRESS, quint16 port = STD_PORT);
    ~Comunicator();

    //PUBLIC METHODS
    void run(QString* username);
    inline QString getUserName() {
        return database.getNomeUtente();
    }

private:
    enum Action {
        ADD = 'A',
        DEL = 'D',
        POS = 'P'
    };

    QHostAddress address;
    quint16 port;
    QTcpSocket *connessione;
    QTextStream stream;

    Database database;

    int cursor_receive;
    int cursor_send;
    QString buffer;
    Action currentAction;
    int active_user_id;

    //PRIVATE METODS
    void svolgiAzione(char buf);
    bool handShaking();
    bool mangiaCarattereDiControllo(char c = '\\');


    
signals:
    void testoAggiunto(int posizione, QString *addString);
    void testoAggiunto(int posizione, QChar addChar);
    void testoRimosso(int posizione, int howmany);
    
public slots:
    void connect_to_host();
    void cambiaPosizione(int pos);
    void aggiungiTesto(int pos, QString *stringa);
    void rimuoviTesto(int pos, int howmany);

private slots:
    void ricevi();
    
};

#endif // COMUNICATOR_H

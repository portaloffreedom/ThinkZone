#include "comunicator.h"
#include <iostream>
using namespace std;

Comunicator::Comunicator(QObject *parent, const QHostAddress &address, quint16 port) :
    QThread(parent)
{
    //TODO memorizzare anche il codice dell'utente
    this->port = port;
    this->address = QHostAddress(address);

    cursor_receive = 0;
    currentAction = ADD;

    //TODO devo fare in modo che le letture dallo "stream" siano bloccanti fino a che non arrivi qualcosa dannazione...
}

Comunicator::~Comunicator()
{
    stream.flush();
    connessione->close();

    delete connessione;
}

bool Comunicator::mangiaCarattereDiControllo(char c)
{
    char c_stream;
    stream>>c_stream;
    if (c_stream != c) {
        cerr<<"non ho trovato il carattere di controllo"<<endl;
        return false;
    }

    return true;
}

bool Comunicator::handShaking()
{
    //TODO stabilire sta stanno parlando la stessa lingua

    //invia nome utente
    QString username = database.getNomeUtente();
    stream<<username<<'\\'<<flush;

    //ricevi id utente
    connessione->waitForReadyRead();

    int id_utente;
    stream>>id_utente;
    if (id_utente <= 0) {
        cerr<<"id_utente non disponibile"<<endl;
        return false;
    }

    database.setIdUtente(id_utente);
    database.setNomeUtente(&username);
    cout<<"l'id dell'utente \""<<username.toStdString()<<"\" è: "<<id_utente<<endl;

    //carattere di controllo
    return mangiaCarattereDiControllo();
}

void Comunicator::connect_to_host()
{
    this->connessione= new QTcpSocket(this);
    connessione->connectToHost(address,port);

    stream.setDevice(connessione);
    if (!handShaking()) {
        cerr<<"ERRORE NELLO STRINGERSI LA MANO: le due mani non sono compatibili"<<endl;
        //TODO gestire l'errore
    }

    connect(connessione,SIGNAL(readyRead()),this,SLOT(ricevi()));
}

void Comunicator::cambiaPosizione(int pos)
{
    if (pos == cursor_send) return;

    stream<<'\\'<<char(POS)<<pos<<'\\'<<flush;
//    cout<<'\\'<<char(POS)<<pos<<'\\'<<flush;
    cursor_send = pos;
}

void Comunicator::aggiungiTesto(int pos, QString *stringa)
{
    cursor_send += stringa->size();
    cambiaPosizione(pos);

    QString buffer(*stringa);
    buffer.replace('\\',"\\\\");

    stream<<buffer<<flush;
//    cout<<stringa->toStdString()<<flush;
}

void Comunicator::rimuoviTesto(int pos, int howmany)
{
    cambiaPosizione(pos+howmany);
    cursor_send -= howmany;

    stream<<'\\'<<char(DEL)<<howmany<<'\\'<<flush;
//    cout<<'\\'<<char(DEL)<<howmany<<'\\'<<flush;
}

void Comunicator::svolgiAzione(char buf)
{
    switch (buf) {
    case '\\':{
        char buf2;
        stream>>buf2;
        switch (buf2) {
        //case 'p':
        case 'P': {
            stream>>cursor_receive;

            if (!mangiaCarattereDiControllo())
                cerr<<"errore nello stream TCP #Delete"<<endl;
            return;
        }
        case 'A':
            currentAction = ADD;
            return;
        //case 'd':
        case 'D': {
            int howmany;
            stream>>howmany;
            cursor_receive -= howmany;
            if (active_user_id != database.getUserId())
                emit testoRimosso(cursor_receive,howmany);

            if (!mangiaCarattereDiControllo())
                cerr<<"errore nello stream TCP #Delete"<<endl;
            return;
        }
        case 'U':
            stream>>active_user_id;
            if (!mangiaCarattereDiControllo())
                cerr<<"errore nello stream TCP #ChangeUser"<<endl;
            return;
//        case 'n':
//            //buf = '\n';
//            this->svolgiAzione('\n');
//            return;
        default:
            cerr<<"errore: azione non disponibile\n"
               "probabilmente adesso tutta la comunicazione è uscita di fasa"
               "e renderà la comunicazione instabile..."<<endl;

        case '\\':
            cout<<"la barra non era un'azione"<<endl;
            //la barra è stata escapata
        }
    }
    default:
        switch (currentAction) {
        case ADD:
//            if (buf == '\n')
//                return;
            if (active_user_id != database.getUserId())
                emit testoAggiunto(cursor_receive,buf);
            cursor_receive++;
            break;
//        case DEL: {
//            stream.seek(stream.pos()-1);
//            int howmany;
//            stream>>howmany;
//            testoRimosso(cursor,howmany);
//            break;
//            }
        default:
            cerr<<"errore"<<endl;
        }
        break;
    }
}

void Comunicator::ricevi()
{
    //disconnect(connessione,SIGNAL(readyRead()),this,SLOT(ricevi()));

    char buf;
    while (!stream.atEnd()) {
        stream>>buf;
        svolgiAzione(buf);
    }

    //connect(connessione,SIGNAL(readyRead()),this,SLOT(ricevi()));
}

void Comunicator::run(QString *username)
{
    this->database.setNomeUtente(username);
    delete username;
    connect_to_host();
    this->setParent(NULL);
//    exec();
}

/*
void Comunicator::ricevi()
{
    char azione;
    int posizione;
    stream>>azione>>posizione;
    stream.flush();
    stream.readLine();
    switch (azione) {
    case opAggiungi: {
        QString addString = stream.readLine();
//        addString->chop(1);

        testoAggiunto(posizione,&addString);
        break;
    }
    case opCancella: {
        int howmanyrem;
        stream>>howmanyrem;
        testoRimosso(posizione,howmanyrem);
        break;
    }
    default:
        cerr<<"che cazzo sta succedendo?"<<endl;
        break;
    }
}
*/

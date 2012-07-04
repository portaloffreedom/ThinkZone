#include "database.h"

Database::Database(QObject *parent) :
    QObject(parent)
{
    this->nome_utente = new QString();
}

Database::~Database()
{
    delete nome_utente;
}

void Database::setNomeUtente(QString *nome_utente)
{
    //spero che questa funzione faccia quello che deve fare...
    //this->nome_utente->swap(*nome_utente);
    this->nome_utente->clear();
    this->nome_utente->append(nome_utente);
}

void Database::setIdUtente(int id)
{
    this->id_utente = id;
}

//inline QString Database::getNomeUtente()
//{
//    return QString(*nome_utente);
//}

int Database::getUserId()
{
    return id_utente;
}

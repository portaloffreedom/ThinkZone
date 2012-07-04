#include "superstring.h"
#include <iostream>
using namespace std;

SuperString::SuperString()
{
    this->total = new QString();
    this->toRebuildTotal = true;

    // crea lista
    testa = new elemListaString(NULL,NULL);
}

SuperString::~SuperString()
{
    delete this->total;

    elemListaString *posAttuale;

    for (int i=1; true; i++) {
        posAttuale = testa->succ;
        delete testa;

        if (posAttuale == NULL)
            return;

        testa = posAttuale;
    }
}

SuperString::elemListaString::elemListaString(elemListaString* prec, elemListaString* succ)
{
    this->elem = new QString;
    this->size = 0;
    this->succ = succ;
    this->prec = prec;
}

SuperString::elemListaString::~elemListaString()
{
    delete this->elem;
}

/**
 * Attenzione! non ricalcola la nuova size per una questione di performance
 *
 * @brief SuperString::elemListaString::sostituisciStringa
 * @param nuova
 */
void SuperString::elemListaString::sostituisciStringa(QString *nuova)
{
    delete this->elem;
    this->elem = nuova;
}

QString *SuperString::getComplete()
{
    if (!toRebuildTotal) {
        return total;
    }
    else {
        delete total; //questa potrebbe essere spostata: nella sezione appena la stringa diventa sporca
        total = new QString();

        // fai l'append di tutta la listas
        //total->append("merda");
        elemListaString *tmp = testa;
        while (tmp != NULL) {
            total->append(tmp->elem);
            tmp = tmp->succ;
        }

        toRebuildTotal = false;
        return total;
    }
}

QString *SuperString::getCompleteWithSeparators()
{

    delete total; //questa potrebbe essere spostata appena perde senso
    total = new QString();

    // fai l'append di tutta la listas
    //total->append("merda");
    elemListaString *tmp = testa;
    while (tmp != NULL) {
        total->append(tmp->elem);
        total->append("[");
        total->append(QString::number(tmp->size));
        total->append("]");
        tmp = tmp->succ;
    }

    toRebuildTotal = true;
    return total;
}

void SuperString::insElem(QString *append, int pos)
{
    if (testa == NULL) {
        testa = new elemListaString(NULL,NULL);
    }

    int dim = append->size();
    if (dim != 0) toRebuildTotal = true;

    if (pos == 0) {
        testa->elem->prepend(append->toAscii());
        testa->size += dim;

        return;
    }

    elemListaString *posAttuale = testa;

    int cont = 0;
    while (posAttuale != NULL && pos >= posAttuale->size) {
        cont++;
        pos -= posAttuale->size;
        posAttuale = posAttuale->succ;
    }

    if (posAttuale == NULL || pos == -1) {
        cout<<"#inserisco in fondo all'ultima stringa"<<endl;
        elemListaString *coda = testa;
        while(coda->succ != NULL) coda = coda->succ;

        coda->elem->append(append);
        coda->size += dim;

        return;
    }

    if (pos == 0) { //non c'è da scindere le stringhe
        cout<<"#inserisco tra 2 stringhe"<<endl;
        posAttuale = posAttuale->prec;
    }
    else { //c'è da scindere le stringhe
        cout<<"#scindo le stringhe. pos="<<pos<<" cont:"<<cont<<endl;
        elemListaString *tmp = new elemListaString(posAttuale,posAttuale->succ);
        posAttuale->succ = tmp;

        QString *vecchio = posAttuale->elem;
        tmp       ->sostituisciStringa(new QString(vecchio->mid(pos)));
        posAttuale->sostituisciStringa(new QString(vecchio->left(pos)));

        tmp->size = posAttuale->size -pos;
        posAttuale->size = pos;
    }

    posAttuale->elem->append(append);
    posAttuale->size += dim;

    return;
}

void SuperString::insElem(const char *append, int pos)
{
    QString appendimi(append);
    insElem(&appendimi, pos);
}

//void SuperString::delElem(int pos, int howmany)
//{
//    if (howmany != 1)
//        //TODO
//        cerr<<"ATTENZIONE!!! "
//              "operazione di eliminazione di più di un elemento "
//              "alla volta ancora non supportata correttamente"<<endl;
//    //rob
//    while((howmany--)>0){
//        cout<<"prov"<<endl;
//        delSingleElem(pos);
//    }
//}

void SuperString::delSingleElem(elemListaString* elemento)
{
    if (elemento->succ != NULL)
        elemento->succ->prec = elemento->prec;

    if (elemento->prec != NULL)
        elemento->prec->succ = elemento->succ;
    else
        testa = elemento->succ;

    if (testa == NULL) {
        testa = new elemListaString(NULL,NULL);
    }

    delete elemento;
}

void SuperString::delElem(int pos, const int howmany){
    if (pos < 0) {
        cerr<<"che cazzo stai cercando di eliminare???"<<endl;
        return;
    }

    this->toRebuildTotal = true;

    elemListaString *posAttuale = testa;

    //int cont = 0;
    while (pos >= posAttuale->size) {
        //cont++;
        pos -= posAttuale->size;
        posAttuale = posAttuale->succ;
    }

    /*
    while(pos+howmany > posAttuale->size) {
        int quanti_eliminare = posAttuale->size-pos+1;
        posAttuale->elem->remove(pos,quanti_eliminare);
        //if (pos == 0) //elimina elemento lista
            //TODO
        //posAttuale->size -= posAttuale->size-pos; //rob: avevi scritto == pos, sto provando a interpretarlo così
        posAttuale->size = pos-1;
        pos = 0;
        howmany -= (quanti_eliminare-1);
        posAttuale = posAttuale->succ;
    }
    posAttuale->elem->remove(pos,howmany);
    posAttuale->size -= howmany;
    */



    //elimina prima stringa
    int quanti_eliminare = howmany;
    if (pos+howmany <= posAttuale->size) {   //elimina solo posizione attuale
        posAttuale->elem->remove(pos,quanti_eliminare);
        posAttuale->size -= quanti_eliminare;
        if (posAttuale->size == 0)
            delSingleElem(posAttuale);
        return;
    }
    quanti_eliminare = posAttuale->size-pos;
    posAttuale->elem->remove(pos,quanti_eliminare);
    posAttuale->size -= (quanti_eliminare);

    if (posAttuale->size == 0) {
        elemListaString* temp = posAttuale->succ;
        delSingleElem(posAttuale);
        posAttuale = temp;
        if (posAttuale == NULL)
            posAttuale = testa;
    }
    else
        posAttuale = posAttuale->succ;

    //elimina stringhe di mezzo
    quanti_eliminare = howmany - quanti_eliminare;

    while (posAttuale->succ != NULL && quanti_eliminare >= posAttuale->size) {
        quanti_eliminare -= posAttuale->size;

        elemListaString* temp = posAttuale->succ;
        delSingleElem(posAttuale);
        posAttuale = temp;
    }

    //elimina ultima stringa

    if (posAttuale->size == quanti_eliminare){
        delSingleElem(posAttuale);
    }
    else {
        posAttuale->elem->remove(0,quanti_eliminare);
        posAttuale->size -= quanti_eliminare;
    }


}

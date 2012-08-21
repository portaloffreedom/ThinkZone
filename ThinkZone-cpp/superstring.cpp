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
//         cout<<"#inserisco in fondo all'ultima stringa"<<endl;
        elemListaString *coda = testa;
        while(coda->succ != NULL) coda = coda->succ;

        coda->elem->append(append);
        coda->size += dim;

        return;
    }

    if (pos == 0) { //non c'è da scindere le stringhe
//         cout<<"#inserisco tra 2 stringhe"<<endl;
        posAttuale = posAttuale->prec;
    }
    else { //c'è da scindere le stringhe
//         cout<<"#scindo le stringhe. pos="<<pos<<" cont:"<<cont<<endl;
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

void errore(QString *s) 
{
  cerr<<"######Errore:|"<<s->toStdString()<<"|\n";
}
void corretto(QString *s)
{
  cout<<"Corretto: "<<s->toStdString()<<endl;
}

bool testString(SuperString *superstringa, char *neutra, char *conSep, int &cont)
{
  cout<<"# "<<cont<<endl;
  QString *stringa = superstringa->getComplete();
  if (stringa != QString(neutra)) {
    errore(stringa);
    return false;
  }
  else 
    corretto(stringa);
  
  stringa = superstringa->getCompleteWithSeparators();
  
  if (stringa != QString(conSep)) {
    errore(stringa);
    return false;
  }
  else 
    corretto(stringa);
  
  cout<<"# "<<cont<<" eseguita"<<endl;
  cont++;
  return true;
}

bool SuperString::Test()
{
  SuperString *prova = new SuperString();
  int i = 0;
  
  //#1
  if (!testString(prova,"","[0]",i))
    return false;

  //#2 ---insert---
  prova->insElem("99", 0);
  if (!testString(prova,"99", "99[2]",i))
    return false;

  //#3
  prova->insElem("12!", 0);
  if (!testString(prova,"12!99", "12!99[5]",i))
    return false;

  //#4
  prova->insElem("###", 3);
  if (!testString(prova,"12!###99", "12!###[6]99[2]",i))
    return false;

  //#5
  prova->insElem("tro", 2);
  if (!testString(prova,"12tro!###99", "12tro[5]!###[4]99[2]",i))
    return false;

  //#6 ----delete----
  prova->delElem(1, 1);
  if (!testString(prova,"1tro!###99", "1tro[4]!###[4]99[2]",i))
    return false;

  //#7
  prova->delElem(4, 1);
  if (!testString(prova,"1tro###99", "1tro[4]###[3]99[2]",i))
    return false;

  //#8
  prova->delElem(4, 3);
  if (!testString(prova,"1tro99", "1tro[4]99[2]",i))
    return false;

  //#9
  prova->delElem(3, 1);
  if (!testString(prova,"1tr99", "1tr[3]99[2]",i))
    return false;

  //#10
  prova->delElem(0, 3);
  if (!testString(prova,"99", "99[2]",i))
    return false;

  //#11
  prova->insElem("porcoDioZoccolo", 0);
  prova->insElem(" ", 8);
  prova->insElem(" ", 5);
  if (!testString(prova,"porco Dio Zoccolo99", "porco [6]Dio [4]Zoccolo99[9]",i))
    return false;

  //#12
  prova->delElem(3, 15);
  if (!testString(prova,"por9", "por[3]9[1]",i))
    return false;
  
  
  return true;
}
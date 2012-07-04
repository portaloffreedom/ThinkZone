#include "inputtesto.h"
#include <iostream>
using namespace std;

inputTesto::inputTesto(QWidget *parent) :
    QTextEdit(parent)
{
    //connect(this->document(),SIGNAL(contentsChange(int,int,int)),this,SLOT(testoCambiato_Slot(int,int,int)));
    //contentsChange ( int position, int charsRemoved, int charsAdded )
    connect(this,SIGNAL(testoAggiunto(int,QString*)),this,SLOT(aggiungiTesto(int,QString*)));
    connect(this,SIGNAL(testoRimosso(int,int)),this,SLOT(rimuoviTesto(int,int)));
}

inputTesto::~inputTesto()
{

}

void inputTesto::setSincTCP(bool attiva)
{
    if (attiva)
        connect(this->document(),SIGNAL(contentsChange(int,int,int)),this,SLOT(testoCambiato_Slot(int,int,int)));
    else
        disconnect(this->document(),SIGNAL(contentsChange(int,int,int)),this,SLOT(testoCambiato_Slot(int,int,int)));

}

void inputTesto::setTextSource(SuperString *testo)
{
    testoSorgente = testo;
}

void inputTesto::textUpdate()
{
    QString merda;
    merda.clear();
    //merda.append(this->testoSorgente->getCompleteWithSeparators());
    merda.append(this->testoSorgente->getComplete());
    this->setText(merda);
}

//parte di invio aggiornamenti:

void inputTesto::testoCambiato_Slot(int position, int charsRemoved, int charsAdded)
{
    cerr<<"pos,rem,add:"<<position<<":"<<charsRemoved<<":"<<charsAdded<<endl;

    if (charsRemoved != 0) {
//        this->testoSorgente->delElem(position, charsRemoved);
        testoRimosso(position,charsRemoved);
    }
    if (charsAdded != 0) {
//        this->testoSorgente->insElem(new QString(this->toPlainText().mid(position,charsAdded)),position);
        QString aggiunta(this->toPlainText().mid(position,charsAdded));
        testoAggiunto(position,&aggiunta);
        position+= charsAdded;
    }

    QTextCursor cursore = this->textCursor();
    cursore.setPosition(position);

    this->setTextCursor(cursore);
    this->testoCambiato();
}


//parte di ricezione aggiornamenti:

void inputTesto::aggiungiTesto(int position, QString *addString)
{
    this->setSincTCP(false);
    cerr<<"pos,add:"<<position<<":"<<addString->toStdString()<<endl;

    this->testoSorgente->insElem(addString,position);
    this->textUpdate();
    this->setSincTCP(true);
}

void inputTesto::aggiungiTesto(int position, QChar addChar)
{
    QString addString(addChar);
    aggiungiTesto(position,&addString);
}

void inputTesto::rimuoviTesto(int position, int howmany)
{
    this->setSincTCP(false);
    cerr<<"pos,rem"<<position<<":"<<howmany<<endl;

    this->testoSorgente->delElem(position,howmany);
    this->textUpdate();
    this->setSincTCP(true);
}

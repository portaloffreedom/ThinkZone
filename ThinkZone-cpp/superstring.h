#ifndef SUPERSTRING_H
#define SUPERSTRING_H

#include <QtGui>

class SuperString
{

public:
    SuperString();
    ~SuperString();

    QString *getComplete();
    QString *getCompleteWithSeparators();

    void insElem ( QString*, int pos=-1 );
    void insElem ( const char *, int pos=-1 );

    void delElem ( int pos, const int howmany=1 );

    static bool Test();

private:

    //dichiarazione classi necessarie interne
    class elemListaString
    {
    public:
        elemListaString ( elemListaString *prec, elemListaString *succ );
        ~elemListaString();
        void sostituisciStringa ( QString *nuova );
        QString *elem;
        int size;
        elemListaString *succ;
        elemListaString *prec;
    };
    //fine dichiarazioni

    elemListaString *testa;
    QString *total;
    bool toRebuildTotal;
    void delSingleElem ( elemListaString* );
};

#endif // SUPERSTRING_H

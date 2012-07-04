#ifndef INPUTTESTO_H
#define INPUTTESTO_H

#include <QTextEdit>
#include "superstring.h"

class inputTesto : public QTextEdit
{
    Q_OBJECT
public:
    explicit inputTesto(QWidget *parent = 0);
    ~inputTesto();
    void setTextSource(SuperString *testo);
    void setSincTCP(bool);

signals:
    void testoCambiato();
//    void cambiaPosizione(int pos);
    void testoAggiunto(int pos, QString *stringa);
    void testoRimosso(int pos, int howmany);
    
public slots:
    void textUpdate();
    void testoCambiato_Slot(int position, int charsRemoved, int charsAdded);
    void aggiungiTesto(int position, QString *addString);
    void aggiungiTesto(int position, QChar addChar);
    void rimuoviTesto(int position, int howmany);

private:
    SuperString *testoSorgente;
    
};

#endif // INPUTTESTO_H

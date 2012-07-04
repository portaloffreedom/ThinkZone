#ifndef DATABASE_H
#define DATABASE_H

#include <QObject>
#include <QString>

class Database : public QObject
{
    Q_OBJECT
public:
    explicit Database(QObject *parent = 0);
    ~Database();

    void setNomeUtente(QString *nome_utente);

    void setIdUtente(int id);

    inline QString getNomeUtente() {
        return QString(*nome_utente);
    }

    int getUserId();
signals:
    
public slots:

private:
    QString *nome_utente;
    int id_utente;

    //TODO lista degli altri utenti come coppia nome - id
    
};

#endif // DATABASE_H

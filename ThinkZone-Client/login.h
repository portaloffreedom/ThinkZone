#ifndef LOGIN_H
#define LOGIN_H

#include <QWidget>
#include <QLineEdit>
#include <QPushButton>

class Login : public QWidget
{
    Q_OBJECT
public:
    explicit Login(QWidget *parent = 0);

private:
    QLineEdit   *usernameform;
    QPushButton *loginButton;
    
signals:
    void login(QString *username);
    
private slots:
    void reactLogin();
    
};

#endif // LOGIN_H

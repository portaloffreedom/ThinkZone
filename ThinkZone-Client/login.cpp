#include "login.h"

#include <QGridLayout>

Login::Login(QWidget *parent) :
    QWidget(parent)
{
    usernameform = new QLineEdit(this);
    loginButton = new QPushButton("Login",this);

    connect(loginButton,SIGNAL(clicked()),this,SLOT(reactLogin()));
    connect(usernameform,SIGNAL(returnPressed()),this,SLOT(reactLogin()));


    QGridLayout *layout = new QGridLayout;
    layout->addWidget(usernameform, 0, 0, Qt::AlignCenter);
    layout->addWidget(loginButton, 0, 1, Qt::AlignCenter);

    this->setLayout(layout);
    this->setWindowTitle("ThinkZone Login");

}

void Login::reactLogin() {

    QString *username = new QString(usernameform->text());
    emit login(username);


    this->hide();
    this->deleteLater();
}

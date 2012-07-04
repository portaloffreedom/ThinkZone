#include <QApplication>
#include "mainwindow.h"
#include "superstring.h"
#include "login.h"
#include <iostream>
using namespace std;

void testSuperString()
{
    cout<<"### inizio_test ###"<<endl;

    SuperString *stringaDiProva = new SuperString();
    stringaDiProva->insElem("123456789");
    cout<<stringaDiProva->getCompleteWithSeparators()->toStdString()<<endl;
    stringaDiProva->insElem("***",3);
    cout<<stringaDiProva->getCompleteWithSeparators()->toStdString()<<endl;
    stringaDiProva->insElem("###",6);
    cout<<stringaDiProva->getCompleteWithSeparators()->toStdString()<<endl;
    stringaDiProva->insElem("porcopio",6);
    cout<<stringaDiProva->getCompleteWithSeparators()->toStdString()<<endl;
    stringaDiProva->insElem("-20-",20);
    cout<<stringaDiProva->getCompleteWithSeparators()->toStdString()<<endl;
    stringaDiProva->insElem("->sono la fine");
    cout<<stringaDiProva->getCompleteWithSeparators()->toStdString()<<endl;

    cout<<endl;
    cout<<stringaDiProva->getComplete()->toStdString()<<endl;
    cout<<stringaDiProva->getComplete()->toStdString()<<endl;

    delete stringaDiProva;

    cout<<"### fine_test ###"<<endl;
}

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    MainWindow w;
    //w.show();
    Login *loginForm = new Login();
    w.connect(loginForm,SIGNAL(login(QString*)),&w,SLOT(StartMainWindow(QString*)));
    loginForm->show();

//    testSuperString();

    return a.exec();
}

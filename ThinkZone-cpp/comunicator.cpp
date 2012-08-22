/*
    <one line to give the program's name and a brief idea of what it does.>
    Copyright (C) 2012  Matteo <email>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

#include <iostream>
using namespace std;
#include "comunicator.h"

Comunicator::Comunicator ( QObject* parent,QHostAddress *serverAddr, quint16 port) : QObject(parent)
{
    this->serverAddr = serverAddr;
    this->port = port;
    
    workerThread = new QThread();
    this->moveToThread(workerThread);
    connect(workerThread,SIGNAL(started()),this,SLOT(process()));
    connect(this,SIGNAL(finished()),workerThread,SLOT(quit()));
    connect(this,SIGNAL(finished()),this,SLOT(deleteLater()));
    connect(workerThread,SIGNAL(finished()),workerThread,SLOT(deleteLater()));
}

Comunicator::~Comunicator()
{
    delete serverAddr;
}


void Comunicator::process()
{
    cout<<"CIAO MONDO DAL THREAD COMUNICATORE!! :D"<<endl;
  
    emit finished();
    return;
}

void Comunicator::runFinished()
{
    cout<<"il thread del comunicatore è terminato\n";
}

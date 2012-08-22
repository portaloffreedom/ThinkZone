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


#ifndef COMUNICATOR_H
#define COMUNICATOR_H

#include <QObject>
#include <QThread>
#include <QHostAddress>

#define STD_PORT 4242
#define STD_ADDRESS QHostAddress::LocalHost

class Comunicator : public QObject
{
  Q_OBJECT
public:
  
    explicit Comunicator(QObject *parent = 0, QHostAddress *address = new QHostAddress(STD_ADDRESS), quint16 port = STD_PORT);
    ~Comunicator();
    
signals:
    void finished();
    void error(QString err);
    
private slots:
    void process();
    
private:
  
  QThread *workerThread;
  quint16 port;
  QHostAddress *serverAddr;
  
public slots:
    void runFinished();
};

#endif // COMUNICATOR_H

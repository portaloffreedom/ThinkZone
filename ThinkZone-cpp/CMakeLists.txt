cmake_minimum_required(VERSION 2.8)

project(thinkzone)

set(QT_USE_QTNETWORK TRUE)

FIND_PACKAGE( Qt4 REQUIRED )

INCLUDE( ${QT_USE_FILE} )

SET ( SOURCES 
  comunicator.cpp 
  thinkzoneapp.cpp
  main.cpp
  superstring.cpp
  )

add_executable(thinkzone ${SOURCES})
 
TARGET_LINK_LIBRARIES( thinkzone ${QT_QTNETWORK_LIBRARY} ${QT_LIBRARIES} )
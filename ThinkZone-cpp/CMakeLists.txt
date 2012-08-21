cmake_minimum_required(VERSION 2.8)

project(thinkzone)

FIND_PACKAGE( Qt4 REQUIRED )

INCLUDE( ${QT_USE_FILE} )

SET ( SOURCES
  mainwindow.cpp
  main.cpp
  superstring.cpp
  )

add_executable(thinkzone ${SOURCES})
 
TARGET_LINK_LIBRARIES( thinkzone ${QT_LIBRARIES} )
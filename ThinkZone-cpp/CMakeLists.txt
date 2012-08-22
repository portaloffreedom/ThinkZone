cmake_minimum_required(VERSION 2.8)

project(thinkzone)

set(QT_USE_QTNETWORK TRUE)

FIND_PACKAGE( Qt4 REQUIRED )

SET ( SOURCES
    mainwindow.cpp 
    comunicator.cpp 
    thinkzoneapp.cpp
    main.cpp
    superstring.cpp
    )
  
SET ( HEADERS_QOBJECT
    comunicator.h
    thinkzoneapp.h
    mainwindow.h
    )
    
SET ( HEADERS
    superstring.h
    ${HEADERS_QOBJECT}
    )

qt4_wrap_cpp( MOCS ${HEADERS_QOBJECT})

INCLUDE( ${QT_USE_FILE} )
add_definitions(${QT_DEFINITIONS})

add_executable(thinkzone ${SOURCES} ${MOCS})

TARGET_LINK_LIBRARIES( thinkzone ${QT_QTNETWORK_LIBRARY} ${QT_LIBRARIES} )
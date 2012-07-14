#!/bin/bash

export GOPATH=`pwd`

go build src/cmd/ThinkZoneServer/ThinkZoneServer.go

# 
# echo
# echo "### installazione del pacchetto network ###"
# cd $GOPATH/src/thinkzone/network/
# go install
# 
# if [ $0 == 0 ]
#   then echo 
#   echo "\n### installazione del pacchetto network ###"
#   cd $GOPATH/src/thinkzone/database/
#   go install
# fi 
# 
# echo
#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

docker-compose -f docker-compose.yml down

docker-compose -f docker-compose.yml up -d


sleep 5

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel create -o orderer.example.com:7050 -c mychannel -f /etc/hyperledger/configtx/channel.tx

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel join -b mychannel.block

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.example.com peer channel update -o orderer.example.com:7050 -c mychannel -f /etc/hyperledger/configtx/Org1MSPanchors.tx

docker exec cli peer chaincode install -n marbles -v 1.0 -p github.com/marbles02/go/
docker exec cli peer chaincode instantiate -o orderer.example.com:7050  -C mychannel -n marbles -v 1.0 -c '{"Args":[""]}'

sleep 5

docker exec cli peer chaincode invoke -C mychannel -n marbles -c '{"Args":["initMarble","marble1","blue","35","tom"]}'
docker exec cli peer chaincode invoke -C mychannel -n marbles -c '{"Args":["initMarble","marble2","red","50","tom"]}'
docker exec cli peer chaincode invoke -C mychannel -n marbles -c '{"Args":["initMarble","marble3","blue","70","tom"]}'

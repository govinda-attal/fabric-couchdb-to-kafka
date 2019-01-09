# Hyperledger Fabric Blockchain world-state data (CouchDB) on Kafka!

![alt text](soln-arch.jpg "Solution Architecture")

## start fabric docker network

> This is a copy from https://github.com/hyperledger/fabric-samples

Start fabric docker network locally
```

cd fabric/basic-network
./start.sh

```

## CoudhDB/Cloudant Kafka Connector

Build connector and copy JARs to connector folder or use the ones in github repo /kafka/connector.

For more details refer to : https://github.com/cloudant-labs/kafka-connect-cloudant

## Setup Kafka

The command below will setup docker containers within same network as above fabric basic network.

```
cd kafka
./setup.sh
sleep 60
./setup-connector.sh
```

## See messages on Kafka Topic
Browse to http://localhost:8000/#/cluster/default/topic/n/marbles/data to view any messages on topic

![alt text](success.png "CouchDB documents on Kafka")


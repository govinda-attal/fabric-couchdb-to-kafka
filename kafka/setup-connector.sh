echo ensure kafka connect is running! sleep for about 60 seconds 

curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"name":"marbles-connector", "config": { "connector.class":"com.ibm.cloudant.kafka.connect.CloudantSourceConnector","tasks.max":"1","cloudant.db.url":"http://couchdb:5984/mychannel_marbles","cloudant.db.username":"admin","cloudant.db.password":"admin","topics":"marbles"} }' \
  http://localhost:8083/connectors

sleep 5

curl -s -X GET http://localhost:8083/connectors/marbles-connector/status
echo ensure kafka connect is running! sleep for about 60 seconds 

curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"name":"couchdb-connector", "config": { "connector.class":"com.ibm.cloudant.kafka.connect.CloudantSourceConnector","tasks.max":"1","cloudant.db.url":"http://couchdb:5984/mychannel_marbles","cloudant.db.username":"admin","cloudant.db.password":"admin","topics":"marbles"} }' \
  http://localhost:8083/connectors

sleep 5

curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"name":"rdb-connector", "config":{"connector.class":"io.confluent.connect.jdbc.JdbcSinkConnector","connection.password":"admin","topics":"out-marble","tasks.max":"1","batch.size":"1","auto.evolve":"true","connection.user":"admin","auto.create":"true","connection.url":"jdbc:postgresql://postgres-db:5432/fabricdb","value.converter":"org.apache.kafka.connect.json.JsonConverter","insert.mode":"upsert","pk.mode":"record_value","pk.fields":"name"}}' \
  http://localhost:8083/connectors

curl -s -X GET http://localhost:8083/connectors/couchdb-connector/status
curl -s -X GET http://localhost:8083/connectors/rdb-connector/status
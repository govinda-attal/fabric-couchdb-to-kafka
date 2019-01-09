docker-compose up -d

sleep 10

docker exec broker kafka-topics --create --topic marbles --partitions 1 --replication-factor 1 --if-not-exists --zookeeper zookeeper:2181

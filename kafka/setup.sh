docker-compose up -d

sleep 10

docker exec broker kafka-topics --create --topic marbles --partitions 1 --replication-factor 1 --if-not-exists --zookeeper zookeeper:2181

docker exec broker kafka-topics --create --topic marbles-grp-table  --partitions 1 --replication-factor 1 --if-not-exists --zookeeper zookeeper:2181 --config cleanup.policy=compact

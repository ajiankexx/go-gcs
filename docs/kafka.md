# Some Usage of kafka
- enter kafka container
```bash
docker exec -it kafka bash
```
create topic:
```bash
kafka-topics \
    --create \
    --bootstrap-server localhost:9092 \
    --replication-factor 1 \
    --partitions 1 \
    --topic email-sender
```

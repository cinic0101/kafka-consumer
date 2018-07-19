# kafka-consumer
A simple kafka consumer, for now support elasticsearch only

## Build:
```bash
docker build -t cinic0101/kafka-consumer .
```

## Docker Run
```bash
docker run cinic0101/kafka-consumer
```

## consumer.yml Template:
```YAML
kafka:
  bootstrap-servers: 192.168.100.1  # Kafka bootstrap servers
  group-id: es  # Kafka group ID
  topics: es  # Kafka topics

destination:
  type: es  # For now, only supports "es", es -> elasticsearch
  params:
    server: 192.168.100.2:9300  # Elasticsearch server
    default-type: default  # Because of elasticsearch 6.x breaking changes, only supports 1 type
```
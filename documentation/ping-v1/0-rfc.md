# v1

`v1` for this application is all about adding stream processing.

## Abstract

This RFC proposes adding stream processing capabilities to the existing ping service using Kafka, Flink, and ScyllaDB. The initial implementation will focus on basic real-time metrics processing.

## Background

The current system stores ping events in PostgreSQL. While this works for basic storage and querying, adding stream processing will enable real-time analytics and metrics.

## High-Level Design

### Current Architecture

```plaintext
Client → Server → PostgreSQL
```

### Proposed Architecture

```plaintext
Client → Server → PostgreSQL
              ↘ Kafka → Flink → ScyllaDB
```

## Components

- **Kafka**: Message broker for ping events
- **Flink**: Stream processing for real-time metrics
- **ScyllaDB**: Storage for processed metrics

## Metrics (Initial Phase)

- **Ping count per client in 5-minute windows**
- **Average ping frequency**
- **Rolling 24-hour activity summary**

## Implementation Phases

### Phase 1: Kafka Integration

#### Setup

- Deploy Kafka using Strimzi operator
- Create `raw-pings` topic
- Configure monitoring (Kafka UI)

#### Server Changes

- Add Kafka producer
- Modify transaction handling
- Add health checks

#### Validation

- Test event production
- Verify transaction atomicity
- Check monitoring

### Phase 2: ScyllaDB Setup

#### Setup

- Deploy ScyllaDB operator
- Create `metrics` keyspace
- Configure monitoring

#### Schema Creation

- Create tables for metrics
- Set up retention policies
- Create necessary indices

#### Validation

- Test connections
- Verify data model
- Check monitoring

### Phase 3: Flink Processing

#### Setup

- Deploy Flink operator
- Configure checkpointing
- Set up monitoring

#### Implementation

- Create processing job
- Configure Kafka source
- Configure ScyllaDB sink

#### Validation

- Test processing pipeline
- Verify metrics accuracy
- Check monitoring

## Detailed Implementation

### Kafka Topic Configuration

```yaml
topics:
    - name: raw-pings
        partitions: 6
        replication-factor: 3
        configs:
            retention.ms: 604800000  # 7 days
            cleanup.policy: delete
```

### ScyllaDB Schema

```sql
CREATE KEYSPACE IF NOT EXISTS metrics
WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': 3};

CREATE TABLE metrics.window_stats (
        client_ip text,
        window_start timestamp,
        window_end timestamp,
        ping_count int,
        avg_interval_ms double,
        PRIMARY KEY ((client_ip), window_start)
) WITH CLUSTERING ORDER BY (window_start DESC);

CREATE TABLE metrics.daily_stats (
        client_ip text,
        day date,
        total_pings counter,
        PRIMARY KEY ((client_ip), day)
) WITH CLUSTERING ORDER BY (day DESC);
```

### Server Implementation

```go
func (s *PingServiceServer) Ping(ctx context.Context, req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error) {
        timestamp := time.Unix(0, req.Msg.TimestampMs*int64(time.Millisecond)).UTC()
        clientIP := getClientIP(ctx)

        // Start transaction
        tx, err := s.db.Begin(ctx)
        if (err != nil) {
                return nil, fmt.Errorf("failed to begin transaction: %v", err)
        }
        defer tx.Rollback(ctx)
        
        // Store in PostgreSQL
        _, err = tx.Exec(ctx, "INSERT INTO pings (pinged_at) VALUES ($1)", timestamp)
        if (err != nil) {
                return nil, fmt.Errorf("failed to store ping: %v", err)
        }
        
        // Prepare Kafka message
        event := PingEvent{
                ClientIP:  clientIP,
                Timestamp: timestamp,
        }
        
        eventBytes, err := json.Marshal(event)
        if (err != nil) {
                return nil, fmt.Errorf("failed to marshal event: %v", err)
        }
        
        // Send to Kafka - if this fails, transaction will be rolled back
        err = s.kafkaWriter.WriteMessages(ctx, kafka.Message{
                Key:   []byte(clientIP),
                Value: eventBytes,
        })
        if (err != nil) {
                return nil, fmt.Errorf("failed to publish event: %v", err)
        }
        
        // Commit transaction only if Kafka write succeeded
        if (err := tx.Commit(ctx); err != nil) {
                return nil, fmt.Errorf("failed to commit transaction: %v", err)
        }
        
        return connect.NewResponse(&pingv1.PingResponse{}), nil
}
```

### Flink Job

```java
public class PingMetricsJob {
        public static void main(String[] args) throws Exception {
                StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

                // Configure Kafka source
                Properties properties = new Properties();
                properties.setProperty("bootstrap.servers", "kafka:9092");
                properties.setProperty("group.id", "ping-metrics");
                
                FlinkKafkaConsumer<PingEvent> kafkaSource = new FlinkKafkaConsumer<>(
                        "raw-pings",
                        new JSONDeserializationSchema<PingEvent>(),
                        properties
                );
                
                // Process 5-minute windows
                DataStream<WindowedMetrics> windowedMetrics = env
                        .addSource(kafkaSource)
                        .keyBy(event -> event.getClientIP())
                        .timeWindow(Time.minutes(5))
                        .aggregate(new PingMetricsAggregator());
                
                // Store in ScyllaDB
                CassandraSink.addSink(windowedMetrics)
                        .setHost("scylla")
                        .setQuery("INSERT INTO metrics.window_stats (client_ip, window_start, window_end, ping_count, avg_interval_ms) VALUES (?, ?, ?, ?, ?)")
                        .build();
                
                env.execute("Ping Metrics");
        }
}
```

### Operational Notes

#### Kafka Topic Inspection

```bash
# List topics
kubectl exec -it kafka-0 -- kafka-topics.sh --list --bootstrap-server localhost:9092

# View topic details
kubectl exec -it kafka-0 -- kafka-topics.sh --describe --topic raw-pings --bootstrap-server localhost:9092

# Monitor consumer group
kubectl exec -it kafka-0 -- kafka-consumer-groups.sh --describe --group ping-metrics --bootstrap-server localhost:9092

# Read messages from topic
kubectl exec -it kafka-0 -- kafka-console-consumer.sh \
        --bootstrap-server localhost:9092 \
        --topic raw-pings \
        --from-beginning \
        --property print.timestamp=true
```

#### Health Checks

- **Kafka**: `kubectl exec kafka-0 -- kafka-broker-api-versions.sh --bootstrap-server localhost:9092`
- **ScyllaDB**: `kubectl exec scylla-0 -- nodetool status`
- **Flink**: Check Flink dashboard (port-forward to 8081)

### Risks and Mitigations

- **Data Loss**: Use proper replication factors and transaction handling
- **Performance**: Monitor and tune Kafka partitions and Flink parallelism
- **Storage**: Implement proper retention policies in both Kafka and ScyllaDB

### Dependencies

- Strimzi Kafka Operator
- Flink Kubernetes Operator
- ScyllaDB Operator
- Prometheus & Grafana for monitoring

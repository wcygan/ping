package com.ping;

import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.api.common.serialization.SimpleStringSchema;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.eventtime.SerializableTimestampAssigner;

import java.time.Duration;

public class PingStreamProcessor {
    public static void main(String[] args) throws Exception {
        final StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

        KafkaSource<String> source = KafkaSource.<String>builder()
            .setBootstrapServers(System.getenv().getOrDefault("KAFKA_BROKERS", "ping-kafka-cluster-kafka-bootstrap.kafka-system:9092"))
            .setTopics("ping-events")
            .setGroupId("ping-processor")
            .setStartingOffsets(OffsetsInitializer.earliest())
            .setValueOnlyDeserializer(new SimpleStringSchema())
            .build();

        // Configure watermark strategy
        WatermarkStrategy<String> watermarkStrategy = WatermarkStrategy
            .<String>forBoundedOutOfOrderness(Duration.ofSeconds(5))
            .withTimestampAssigner((event, timestamp) -> System.currentTimeMillis());

        // Read from Kafka with watermarks
        env.fromSource(source, watermarkStrategy, "Kafka Source")
           .print(); // Just print to stdout for now

        env.execute("Ping Stream Processor");
    }
}

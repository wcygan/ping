package com.ping;

import build.buf.gen.ping.v1.PingRequest;
import com.esotericsoftware.minlog.Log;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.api.common.serialization.AbstractDeserializationSchema;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.streaming.api.CheckpointingMode;

import java.time.Duration;
import java.io.IOException;

public class PingStreamProcessor {
    private static class PingRequestDeserializer extends AbstractDeserializationSchema<PingRequest> {
        @Override
        public PingRequest deserialize(byte[] message) throws IOException {
            return PingRequest.parseFrom(message);
        }
    }

    public static void main(String[] args) throws Exception {
        final StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

        // Enable checkpointing for exactly-once processing
        env.enableCheckpointing(Duration.ofSeconds(10).toMillis());
        env.getCheckpointConfig().setCheckpointingMode(CheckpointingMode.EXACTLY_ONCE);
        env.getCheckpointConfig().setMinPauseBetweenCheckpoints(Duration.ofSeconds(5).toMillis());
        env.getCheckpointConfig().setCheckpointTimeout(Duration.ofMinutes(1).toMillis());
        env.getCheckpointConfig().setMaxConcurrentCheckpoints(1);

        Log.info("Starting PingStreamProcessor");

        KafkaSource<PingRequest> source = KafkaSource.<PingRequest>builder()
            .setBootstrapServers(System.getenv().getOrDefault("KAFKA_BROKERS", "ping-kafka-cluster-kafka-bootstrap.kafka-system:9092"))
            .setTopics("ping-events")
            .setGroupId("ping-processor")
            .setStartingOffsets(OffsetsInitializer.earliest())
            .setValueOnlyDeserializer(new PingRequestDeserializer())
            .setProperty("value.deserializer", "org.apache.kafka.common.serialization.ByteArrayDeserializer")
            .build();

        // Configure watermark strategy using the timestamp from PingRequest
        WatermarkStrategy<PingRequest> watermarkStrategy = WatermarkStrategy
            .<PingRequest>forBoundedOutOfOrderness(Duration.ofSeconds(5))
            .withTimestampAssigner((event, timestamp) -> event.getTimestampMs());

        // Read from Kafka with watermarks
        String redisHost = System.getenv().getOrDefault("REDIS_HOST", "ping-cache");
        int redisPort = Integer.parseInt(System.getenv().getOrDefault("REDIS_PORT", "6379"));

        env.fromSource(source, watermarkStrategy, "Kafka Source")
           .addSink(new BufferingRedisSink(redisHost, redisPort))
           .name("Redis Sink");

        env.execute("Ping Counter");
    }
}

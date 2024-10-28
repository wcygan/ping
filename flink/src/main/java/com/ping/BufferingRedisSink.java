package com.ping;

import build.buf.gen.ping.v1.PingRequest;
import org.apache.flink.api.common.state.ListState;
import org.apache.flink.api.common.state.ListStateDescriptor;
import org.apache.flink.runtime.state.FunctionInitializationContext;
import org.apache.flink.runtime.state.FunctionSnapshotContext;
import org.apache.flink.streaming.api.checkpoint.CheckpointedFunction;
import org.apache.flink.streaming.api.functions.sink.RichSinkFunction;
import org.apache.flink.configuration.Configuration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPool;
import redis.clients.jedis.JedisPoolConfig;
import redis.clients.jedis.Pipeline;
import redis.clients.jedis.exceptions.JedisException;

import java.io.IOException;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;

public class BufferingRedisSink extends RichSinkFunction<PingRequest> 
        implements CheckpointedFunction {
    
    private static final Logger LOG = LoggerFactory.getLogger(BufferingRedisSink.class);
    private static final int BATCH_SIZE = 1000;
    private static final Duration BATCH_TIMEOUT = Duration.ofSeconds(1);
    private static final int MAX_RETRY_ATTEMPTS = 3;
    private static final Duration RETRY_DELAY = Duration.ofMillis(100);
    
    private final String redisHost;
    private final int redisPort;
    
    private transient JedisPool jedisPool;
    private transient List<PingRequest> buffer;
    private transient ListState<PingRequest> checkpointedState;
    private transient long lastFlushTime;

    public BufferingRedisSink(String redisHost, int redisPort) {
        this.redisHost = redisHost;
        this.redisPort = redisPort;
    }

    @Override
    public void open(Configuration parameters) {
        JedisPoolConfig poolConfig = new JedisPoolConfig();
        poolConfig.setMaxTotal(8);
        poolConfig.setMaxIdle(4);
        poolConfig.setMinIdle(1);
        poolConfig.setTestOnBorrow(true);
        poolConfig.setTestWhileIdle(true);
        poolConfig.setMaxWait(Duration.ofSeconds(2));
        
        jedisPool = new JedisPool(poolConfig, redisHost, redisPort);
        buffer = new ArrayList<>();
        lastFlushTime = System.currentTimeMillis();
        
        LOG.info("Initialized BufferingRedisSink with host={}, port={}", redisHost, redisPort);
    }

    @Override
    public void invoke(PingRequest value, Context context) throws IOException {
        LOG.info("Received ping: {}", value);

        buffer.add(value);

        if (shouldFlush()) {
            flush();
        }
    }

    private boolean shouldFlush() {
        return buffer.size() >= BATCH_SIZE || 
               System.currentTimeMillis() - lastFlushTime >= BATCH_TIMEOUT.toMillis();
    }

    private void flush() throws IOException {
        if (buffer.isEmpty()) {
            return;
        }

        int attempts = 0;
        boolean success = false;
        Exception lastException = null;

        while (attempts < MAX_RETRY_ATTEMPTS && !success) {
            LOG.info("Flushing {} records to Redis (attempt {}/{})", buffer.size(), attempts + 1, MAX_RETRY_ATTEMPTS);

            try (Jedis jedis = jedisPool.getResource()) {
                Pipeline pipeline = jedis.pipelined();
                
                for (PingRequest ping : buffer) {
                    // Increment total counter
                    pipeline.hincrBy("ping:counters", "total", 1);
                    
                    // Increment per-minute counter
                    long minuteTimestamp = ping.getTimestampMs() / (60 * 1000);
                    pipeline.hincrBy("ping:counters", "minute:" + minuteTimestamp, 1);
                    
                    // Keep a sliding window of the last 60 minutes
                    pipeline.expire("ping:counters", 3600); // 1 hour TTL
                }
                
                pipeline.sync();
                success = true;
                
                LOG.debug("Successfully flushed {} records to Redis", buffer.size());
                
            } catch (JedisException e) {
                lastException = e;
                attempts++;
                
                if (attempts < MAX_RETRY_ATTEMPTS) {
                    LOG.warn("Failed to flush to Redis (attempt {}/{}), retrying...", 
                            attempts, MAX_RETRY_ATTEMPTS, e);
                    try {
                        Thread.sleep(RETRY_DELAY.toMillis());
                    } catch (InterruptedException ie) {
                        Thread.currentThread().interrupt();
                        throw new IOException("Interrupted during retry delay", ie);
                    }
                }
            }
        }

        if (!success) {
            LOG.error("Failed to flush to Redis after {} attempts", MAX_RETRY_ATTEMPTS);
            throw new IOException("Failed to flush to Redis", lastException);
        }

        buffer.clear();
        lastFlushTime = System.currentTimeMillis();
    }

    @Override
    public void snapshotState(FunctionSnapshotContext context) throws Exception {
        LOG.debug("Taking snapshot of state, {} items in buffer", buffer.size());
        checkpointedState.clear();
        
        // First flush any buffered elements
        flush();
        
        // Checkpoint any remaining items in the buffer
        for (PingRequest ping : buffer) {
            checkpointedState.add(ping);
        }
        
        LOG.debug("State snapshot completed, {} items in buffer", buffer.size());
    }

    @Override
    public void initializeState(FunctionInitializationContext context) throws Exception {
        LOG.info("Initializing state from checkpoint");
        ListStateDescriptor<PingRequest> descriptor =
            new ListStateDescriptor<>("buffered-pings", PingRequest.class);
        checkpointedState = context.getOperatorStateStore().getListState(descriptor);

        if (context.isRestored()) {
            // Restore state after failure
            for (PingRequest ping : checkpointedState.get()) {
                buffer.add(ping);
            }
            LOG.info("Restored {} items from checkpoint", buffer.size());
        }
    }

    @Override
    public void close() throws Exception {
        LOG.info("Closing BufferingRedisSink");
        if (buffer != null && !buffer.isEmpty()) {
            try {
                flush();
            } catch (Exception e) {
                LOG.error("Error flushing buffer during close", e);
            }
        }
        
        if (jedisPool != null) {
            jedisPool.close();
            LOG.info("Closed Redis connection pool");
        }
    }
}

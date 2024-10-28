# Build stage
FROM --platform=linux/arm64 maven:3.8.4-openjdk-11 AS builder

WORKDIR /build
COPY pom.xml .
COPY src ./src

# Build the application
RUN mvn clean package

# Final stage
FROM --platform=linux/arm64 apache/flink:1.18.1

# Copy the application jar from builder stage
COPY --from=builder /build/target/ping-processor-1.0-SNAPSHOT.jar /opt/flink/usrlib/ping-processor-1.0-SNAPSHOT.jar

# Set the entrypoint to the Flink job
ENV FLINK_HOME=/opt/flink
ENV PATH="${FLINK_HOME}/bin:${PATH}"

# Set security related options
ENV JAVA_SECURITY_PROPERTIES="-Djava.security.egd=file:/dev/./urandom"

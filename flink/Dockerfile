# Build stage
FROM maven:3.8.4-openjdk-11 AS builder

WORKDIR /build

# Copy the pom.xml file and download dependencies
COPY pom.xml .
RUN mvn dependency:go-offline

# Copy the source code and build the application
COPY src ./src
RUN mvn clean package

# Final stage
FROM apache/flink:1.18.1

# Copy the application jar from builder stage
COPY --from=builder /build/target/ping-processor-1.0-SNAPSHOT.jar /opt/flink/usrlib/ping-processor-1.0-SNAPSHOT.jar

# Set the entrypoint to the Flink job
ENV FLINK_HOME=/opt/flink
ENV PATH="${FLINK_HOME}/bin:${PATH}"

# Set security related options
ENV JAVA_SECURITY_PROPERTIES="-Djava.security.egd=file:/dev/./urandom"
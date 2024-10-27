FROM maven:3.8-openjdk-11 AS builder
WORKDIR /app
COPY pom.xml .
COPY src src
RUN mvn clean package

FROM flink:1.17.1
COPY --from=builder /app/target/ping-processor-1.0-SNAPSHOT.jar /opt/flink/usrlib/ping-processor.jar
FROM maven:3.8-openjdk-11 AS builder
WORKDIR /app
COPY pom.xml .
COPY src src
RUN mvn clean package

FROM flink:1.17.1
COPY --from=builder /app/target/ping-processor-1.0-SNAPSHOT.jar /opt/flink/usrlib/ping-processor.jar

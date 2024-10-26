package config

import "os"

type Config struct {
    DBHost       string
    DBUser       string
    DBPassword   string
    DBName       string
    KafkaBrokers []string
}

func LoadConfig() Config {
    return Config{
        DBHost:       getEnvOrDefault("DB_HOST", "localhost"),
        DBUser:       getEnvOrDefault("DB_USER", "pinguser"),
        DBPassword:   getEnvOrDefault("DB_PASSWORD", "ping123"),
        DBName:       getEnvOrDefault("DB_NAME", "pingdb"),
        KafkaBrokers: []string{"ping-kafka-cluster-kafka-bootstrap.kafka-system.svc:9092"},
    }
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

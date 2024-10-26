package errors

import "fmt"

type StorageError struct {
    Err error
}

func (e *StorageError) Error() string {
    return fmt.Sprintf("storage error: %v", e.Err)
}

type KafkaError struct {
    Err error
}

func (e *KafkaError) Error() string {
    return fmt.Sprintf("kafka error: %v", e.Err)
}

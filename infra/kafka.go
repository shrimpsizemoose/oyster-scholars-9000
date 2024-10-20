package infra

import (
	"context"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/shrimpsizemoose/trekker/logger"
)

type KafkaConfig struct {
	Addr  string
	Topic string
}

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter(cfg KafkaConfig) *KafkaWriter {
	return &KafkaWriter{
		writer: &kafka.Writer{
			Addr:        kafka.TCP(cfg.Addr),
			Topic:       cfg.Topic,
			Balancer:    &kafka.LeastBytes{},
			ErrorLogger: logger.Error,
		},
	}
}

func (w *KafkaWriter) WriteMessages(ctx context.Context, messages [][]byte) error {
	kafkaMessages := make([]kafka.Message, len(messages))
	for i, msg := range messages {
		kafkaMessages[i] = kafka.Message{Value: msg}
	}

	err := w.writer.WriteMessages(ctx, kafkaMessages...)
	if err != nil {
		return fmt.Errorf("не удалось записать сообщения в Kafka: %w", err)
	}

	return nil
}

func (w *KafkaWriter) Close() error {
	return w.writer.Close()
}

type KafkaReader struct {
	reader *kafka.Reader
}

func NewKafkaReader(cfg KafkaConfig) *KafkaReader {
	return &KafkaReader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.Addr},
			Topic:   cfg.Topic,
		}),
	}
}

func (r *KafkaReader) ReadMessage(ctx context.Context) ([]byte, error) {
	msg, err := r.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("Kafka: не удалось прочитать: %w", err)
	}
	return msg.Value, nil
}

func (r *KafkaReader) Close() error {
	return r.reader.Close()
}

func CheckKafkaTopic(addr, topic string) (bool, error) {
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		return false, fmt.Errorf("[Kafka] не удалось подключиться: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return false, fmt.Errorf("[Kafka] не удалось прочитать партиции: %w", err)
	}

	for _, p := range partitions {
		if p.Topic == topic {
			return true, nil
		}
	}

	return false, nil
}

func WaitForKafka(addr string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("не дождался связи с Kafka")
		default:
			conn, err := kafka.Dial("tcp", addr)
			if err == nil {
				conn.Close()
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}
}

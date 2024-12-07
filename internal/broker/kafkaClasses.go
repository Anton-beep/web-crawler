package broker

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
	"web-crawler/internal/config"
)

// 💀 💀 💀
type SitesKafka struct {
	cfg      *config.KafkaConfig
	producer *kafka.Conn
	consumer *kafka.Reader
}

type Message struct {
	Link      string `json:"link"`
	ProjectId string `json:"projectId"`
	Depth     int    `json:"depth"`
}

func New(cfg *config.Config, createProducer bool, createConsumer bool) *SitesKafka {
	var conn *kafka.Conn
	var r *kafka.Reader
	if createProducer {
		conn, _ = kafka.DialLeader(context.Background(), "tcp", cfg.Kafka.Address, cfg.Kafka.Topic, cfg.Kafka.Partition)
	}
	if createConsumer {
		r = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{cfg.Kafka.Address},
			GroupID:  cfg.Kafka.SitesGroupID,
			Topic:    cfg.Kafka.Topic,
			MaxBytes: 10e6, // 10MB
		})
	}
	return &SitesKafka{cfg: &cfg.Kafka, producer: conn, consumer: r}
}

func (s SitesKafka) AddSiteToParse(link string, projectId string, depth int) error {
	// to produce message

	message := Message{
		Link:      link,
		ProjectId: projectId,
		Depth:     depth,
	}
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = s.producer.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return err
	}
	_, err = s.producer.WriteMessages(
		kafka.Message{Value: msg},
	)
	zap.S().Debug("Message sent to partition")
	if err != nil {
		return err
	}

	return nil
}

func (s SitesKafka) CheckSitesToParse() (*Message, error) {
	// make a new reader that consumes from topic-A
	m, err := s.consumer.ReadMessage(context.Background())
	if err != nil {
		return &Message{}, err
	}
	zap.S().Debug("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	toReturn := Message{}
	err = json.Unmarshal(m.Value, &toReturn)
	if err != nil {
		return &Message{}, err
	}
	return &toReturn, nil
}

func (s SitesKafka) Close() error {
	if err := s.producer.Close(); err != nil {
		return err
	}
	if err := s.consumer.Close(); err != nil {
		return err
	}
	return nil
}

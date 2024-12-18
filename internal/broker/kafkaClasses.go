package broker

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
	"web-crawler/internal/config"
	"web-crawler/internal/utils"
)

// Kafka is a struct that contains the kafka connection and the kafka config
type Kafka struct {
	cfg       *config.KafkaConfig
	producer  *kafka.Conn
	consumer  *kafka.Reader
	kafkaType string
}

// Message is a struct that contains the link to parse, the project id and the depth
type Message struct {
	Link        string `json:"link"`
	ProjectId   string `json:"projectId"`
	Depth       int    `json:"depth"`
	AnalyseType string `json:"analyseType"`
}

func getProducer(kafkaType string, cfg *config.KafkaConfig) (*kafka.Conn, error) {
	var conn *kafka.Conn
	var err error
	if kafkaType == "sites" {
		conn, err = kafka.DialLeader(context.Background(), "tcp", cfg.Address, cfg.SitesTopic, cfg.Partition)
	} else {
		conn, err = kafka.DialLeader(context.Background(), "tcp", cfg.Address, cfg.AnalyseTopic, cfg.Partition)
	}
	return conn, err
}

// New is a function that creates a new SitesKafka struct
func New(cfg *config.Config, createProducer bool, createConsumer bool, kafkaType string) *Kafka {
	var conn *kafka.Conn
	var r *kafka.Reader
	var err error

	if createProducer {
		conn, err = getProducer(kafkaType, &cfg.Kafka)
		if err != nil {
			zap.S().Error("Error while creating kafka producer")
			return nil
		}
	}

	if createConsumer {
		if kafkaType == "sites" {
			r = kafka.NewReader(kafka.ReaderConfig{
				Brokers:  []string{cfg.Kafka.Address},
				GroupID:  cfg.Kafka.SitesGroupID,
				Topic:    cfg.Kafka.SitesTopic,
				MaxBytes: 10e6, // 10MB
			})
		} else {
			r = kafka.NewReader(kafka.ReaderConfig{
				Brokers:  []string{cfg.Kafka.Address},
				GroupID:  cfg.Kafka.AnalyseGroupID,
				Topic:    cfg.Kafka.AnalyseTopic,
				MaxBytes: 10e6, // 10MB
			})
		}
	}
	return &Kafka{cfg: &cfg.Kafka, producer: conn, consumer: r, kafkaType: kafkaType}
}

func (s *Kafka) writeMessages(message Message) error {
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

	return err
}

func (s *Kafka) produce(message Message) error {
	return utils.RetryCount(4, 1*time.Second, nil, func() error {
		err := s.writeMessages(message)

		if err != nil {
			zap.S().Error("Broken pipe error, retrying...")
			err = s.producer.Close()
			if err != nil {
				zap.S().Error("Error while closing producer: ", err)
			}

			var conn *kafka.Conn
			conn, err = getProducer(s.kafkaType, s.cfg)
			if err != nil {
				return err
			}

			s.producer = conn
			zap.S().Debug("Producer reconnected")

			err = s.writeMessages(message)
		}
		zap.S().Debug("Message sent to partition")

		return err
	})
}

func (s *Kafka) ProduceAnalyse(projectId string, analyseType string) error {
	message := Message{
		Link:        "",
		ProjectId:   projectId,
		Depth:       0,
		AnalyseType: analyseType,
	}
	return s.produce(message)
}

// ProduceSite is a function that adds a site to parse to the kafka topic
func (s *Kafka) ProduceSite(link string, projectId string, depth int) error {
	message := Message{
		Link:        link,
		ProjectId:   projectId,
		Depth:       depth,
		AnalyseType: "",
	}
	return s.produce(message)
}

// Consume is a function that reads a message from the kafka topic
func (s *Kafka) Consume() (*Message, error) {
	// make a new reader that consumes from topic-A
	m, err := s.consumer.ReadMessage(context.Background())
	if err != nil {
		return &Message{}, err
	}
	zap.S().Debugf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	toReturn := Message{}
	err = json.Unmarshal(m.Value, &toReturn)
	if err != nil {
		return &Message{}, err
	}
	return &toReturn, nil
}

// Close is a function that closes the kafka connection
func (s *Kafka) Close() error {
	if err := s.producer.Close(); err != nil {
		return err
	}
	if err := s.consumer.Close(); err != nil {
		return err
	}
	return nil
}

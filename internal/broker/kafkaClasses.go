package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
	"web-crawler/internal/config"
)

// 💀 💀 💀
type SitesKafka struct {
	cfg *config.KafkaConfig
}

type Message struct {
	Link      string `json:"link"`
	ProjectId string `json:"projectId"`
	Depth     int    `json:"depth"`
}

func New(cfg *config.Config) *SitesKafka {
	return &SitesKafka{cfg: &cfg.Kafka}
}

func (s SitesKafka) AddSiteToParse(link string, projectId string, depth int) error {
	// to produce messages

	conn, err := kafka.DialLeader(context.Background(), "tcp", s.cfg.Address, s.cfg.Topic, s.cfg.Partition)
	if err != nil {
		return err
	}

	message := Message{
		Link:      link,
		ProjectId: projectId,
		Depth:     depth,
	}
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return err
	}
	_, err = conn.WriteMessages(
		kafka.Message{Value: msg},
	)
	if err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}
	return nil
}

func (s SitesKafka) CheckSitesToParse() (*Message, error) {
	// make a new reader that consumes from topic-A
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{s.cfg.Address},
		GroupID:  s.cfg.SitesGroupID,
		Topic:    s.cfg.Topic,
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()
	m, err := r.ReadMessage(context.Background())
	if err != nil {
		return &Message{}, err
	}
	//TODO debugging fmt.Print is not cool. Use zap.S().Debug()
	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	toReturn := Message{}
	err = json.Unmarshal(m.Value, &toReturn)
	if err != nil {
		return &Message{}, err
	}
	return &toReturn, nil
}

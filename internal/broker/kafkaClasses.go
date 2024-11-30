package broker

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
)

type Message struct {
	link      string `json:"link"`
	projectId int    `json:"projectId"`
}

type SitesKafka struct {
	consumer *kafka.Consumer
	producer *kafka.Producer
	topic    string
}

func (s *SitesKafka) AddSiteToParse(link string, projectId int) error {
	message := Message{
		link:      link,
		projectId: projectId,
	}
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}
	deliveryChan := make(chan kafka.Event)
	err = s.producer.Produce(kafkaMessage, deliveryChan)
	if err != nil {
		return err
	}
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	close(deliveryChan)

	return nil
}

func (s *SitesKafka) CheckSitesToParse(link string, projectId int) error {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Received signal %v: terminating\n", sig)
			run = false
		default:
			// Poll for Kafka messages
			ev := s.consumer.Poll(1)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				// Process the consumed message
				fmt.Printf("Received message from topic %s: %s\n", *e.TopicPartition.Topic, string(e.Value))
			case kafka.Error:
				// Handle Kafka errors
				fmt.Printf("Error: %v\n", e)
			}
		}
	}
	return nil
}

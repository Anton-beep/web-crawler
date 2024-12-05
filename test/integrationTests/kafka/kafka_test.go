package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"os"
	"sync"
	"testing"
	"time"
	"web-crawler/internal/broker"
	"web-crawler/internal/config"
)

var (
	cfg *config.Config
)

func TestMain(m *testing.M) {
	cfg = config.NewConfig("../../../configs/.env")
	config.InitLogger(true)
	zap.S().Debug(cfg)
	code := m.Run()
	os.Exit(code)
}

func TestConnectionToKafka(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}
	zap.S().Debug("Testing connection")
	// to produce messages
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		zap.S().Error("Failed to connect to kafka dealer", err)
	}

	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		zap.S().Error("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		zap.S().Error("failed to close writer:", err)
	}

	conn1, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		zap.S().Error("failed to dial leader:", err)
	}

	_ = conn1.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn1.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	if err := batch.Close(); err != nil {
		zap.S().Error("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		zap.S().Error("failed to close connection:", err)
	}

	zap.S().Debug("Testing connection finished")
}

func TestKafkaSingle(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}
	siteConf := *broker.New(cfg)
	err := siteConf.AddSiteToParse("link", "1", 0)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	msg, err := siteConf.CheckSitesToParse()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if msg.Link != "link" || msg.ProjectId != "1" {
		t.Errorf("Error: %v", msg)
	}
}

func TestKafkaMultiple(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}
	siteConf := *broker.New(cfg)
	for j := 0; j < 5; j++ {
		err := siteConf.AddSiteToParse(fmt.Sprintf("link %v", j), fmt.Sprintf("%v", j), 0)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	}
	var (
		wg     sync.WaitGroup
		answer = make([]broker.Message, 0)
		mu     sync.Mutex
	)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			msg, err := siteConf.CheckSitesToParse()
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			mu.Lock()
			answer = append(answer, *msg)
			mu.Unlock()
		}()
	}
	wg.Wait()
	time.Sleep(10 * time.Second)
}

//func TestLoginUser(t *testing.T) {
//	err := LoginUser(user)
//	if err != nil {
//		t.Errorf("Error: %v", err)
//	}
//}
//
//func TestRegisterUser(t *testing.T) {
//	err := RegisterUser(user)
//	if err != nil {
//		t.Errorf("Error: %v", err)
//	}
//}
//
//func TestIsFinished(t *testing.T) {
//	err := IsFinished(user)
//	if err != nil {
//		t.Errorf("Error: %v", err)
//	}
//}
//
//func TestGetLoginRequest(t *testing.T) {
//	err := GetLoginRequest()
//	if err != nil {
//		t.Errorf("Error: %v", err)
//	}
//}
//
//func TestGetRegisterRequest(t *testing.T) {
//	err := GetRegisterRequest()
//	if err != nil {
//		t.Errorf("Error: %v", err)
//	}
//}
//
//func TestFinishedLoginOrRegister(t *testing.T) {
//	err := FinishedLoginOrRegister(user, err)
//	if err != nil {
//		t.Errorf("Error: %v", err)
//	}
//}

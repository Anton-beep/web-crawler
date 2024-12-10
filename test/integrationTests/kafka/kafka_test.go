package kafka

import (
	"fmt"
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
	code := m.Run()
	os.Exit(code)
}

func TestKafkaSingle(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}
	siteConf := *broker.New(cfg, true, true, "sites")
	defer siteConf.Close()
	err := siteConf.ProduceSite("link", "1", 0)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	_, err = siteConf.Consume()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestKafkaMultiple(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}
	siteConf := *broker.New(cfg, true, false, "sites")
	//defer siteConf.Close()
	for j := 0; j < 5; j++ {
		err := siteConf.ProduceSite(fmt.Sprintf("link %v", j), fmt.Sprintf("%v", j), 0)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	}
	var (
		wg sync.WaitGroup
	)
	for i := 6; i < 11; i++ {
		wg.Add(1)
		go func() {
			siteConf := *broker.New(cfg, true, true, "sites")
			defer func() {
				siteConf.Close()
				wg.Done()
			}()
			err := siteConf.ProduceSite(fmt.Sprintf("link %v", i), fmt.Sprintf("%v", i), 0)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			_, err = siteConf.Consume()
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			_, err = siteConf.Consume()
			if err != nil {
				t.Errorf("Error: %v", err)
			}
		}()
	}
	wg.Wait()
	time.Sleep(10 * time.Second)
}

func TestKafkaAnalyseSingle(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}
	siteConf := *broker.New(cfg, true, true, "analyse")
	defer siteConf.Close()
	err := siteConf.ProduceAnalyse("link", "1", 0, "type")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	msg, err := siteConf.Consume()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if msg.Link != "link" || msg.ProjectId != "1" {
		t.Errorf("Error: %v", msg)
	}
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

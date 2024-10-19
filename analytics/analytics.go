package analytics

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shrimpsizemoose/trekker/logger"
)

type Tracker interface {
	Ping(eventType string, additionalData map[string]string)
	PingStart()
	PingFinish()
}

type Analytics struct {
	BaseURL       string
	SkipTLS       bool
	CommonData    map[string]string
	SecretHeaders map[string]string
}

func NewAnalytics(baseURL string, skipTLS bool, commonData, secretHeaders map[string]string) Tracker {
	return &Analytics{
		BaseURL:       baseURL,
		SkipTLS:       skipTLS,
		CommonData:    commonData,
		SecretHeaders: secretHeaders,
	}
}

func (a *Analytics) sendEvent(eventType string, additionalData map[string]string) error {
	data := make(map[string]string)
	for k, v := range a.CommonData {
		data[k] = v
	}
	for k, v := range additionalData {
		data[k] = v
	}
	data["event_type"] = eventType
	data["local_datetime"] = time.Now().String()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	http.DefaultClient.Timeout = 2 * time.Second
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: a.SkipTLS}

	req, err := http.NewRequest("POST", a.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	for k, v := range a.SecretHeaders {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error.Fatalf("Что-то не так с аналитикой: я не смог тебя посчитать. Это плохо. Напиши координатору пжлст и приложи скриншот. Спасибо 🐳")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		logger.Error.Fatalf("Неправильное сочетание студента-токена, перепроверь что всё вводишь правильно. Ожидал статус 200 OK, получил - %s", resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error.Fatalf("Ой. Я пытался тебя посчитать, но не смог убедиться что всё ок. Напиши куратору - это важный кейс. Ожидал статус 200 OK, получил - %s", resp.Status)
	}

	return nil
}

func (a *Analytics) Ping(eventType string, additionalData map[string]string) {
	err := a.sendEvent(eventType, additionalData)
	if err != nil {
		logger.Error.Fatalf("Failed to send analytics: %v", err)
	}
}

func (a *Analytics) PingStart() {
	a.Ping("000_lab_start", nil)
}

func (a *Analytics) PingFinish() {
	a.Ping("100_lab_finish", nil)
}

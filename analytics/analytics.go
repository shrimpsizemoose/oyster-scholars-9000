package analytics

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/shrimpsizemoose/trekker/logger"
)

type Tracker interface {
	Ping(eventType string, additionalData map[string]string)
	PingStart()
	PingFinish()
}

type Config struct {
	BaseURL       string
	SkipTLS       bool
	CommonData    map[string]string
	SecretHeaders map[string]string
}

type Analytics struct {
	config  Config
	verbose bool
}

func NewAnalytics(config Config) Tracker {
	_, verbose := os.LookupEnv("TREKKER_VERBOSE")
	return &Analytics{
		config:  config,
		verbose: verbose,
	}
}

func (a *Analytics) sendEvent(eventType string, additionalData map[string]string) error {
	data := make(map[string]string)
	for k, v := range a.config.CommonData {
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

	http.DefaultClient.Timeout = 3 * time.Second
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: a.config.SkipTLS}

	req, err := http.NewRequest("POST", a.config.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	for k, v := range a.config.SecretHeaders {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		if a.verbose {
			logger.Error.Println(err)
		}
		return fmt.Errorf("Что-то не так с аналитикой: я не смог тебя посчитать. Надо проверить сеть, а если не поможет -- напиши координатору пжлст и приложи скриншот. Спасибо 🐳.")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		if a.verbose {
			logger.Error.Println(resp)
		}
		return fmt.Errorf("Неправильное сочетание студента-токена, перепроверь что всё вводишь правильно. Ожидал статус 200 OK, получил - %s", resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		if a.verbose {
			logger.Error.Println(resp)
		}
		return fmt.Errorf("Ой. Я пытался тебя посчитать, но не смог убедиться что всё ок. Надо проверить сеть, а если не поможет -- напиши координатору пжлст и приложи скриншот. Я ожидал статус 200 OK, получил - %s", resp.Status)
	}

	return nil
}

func (a *Analytics) Ping(eventType string, additionalData map[string]string) {
	err := a.sendEvent(eventType, additionalData)
	if a.verbose {
		logger.Warn.Printf("Аналитика: event %s", eventType)
	}
	if err != nil {
		if !a.verbose {
			logger.Warn.Println("Чтобы получить чуть больше информации об ошибке, запусти меня ещё раз с TREKKER_VERBOSE=da")
		}
		logger.Error.Fatalf("Failed to send analytics: %v", err)
	}
}

func (a *Analytics) PingStart() {
	if a.verbose {
		logger.Warn.Println("Аналитика стартует")
	}
	a.Ping("000_lab_start", nil)
}

func (a *Analytics) PingFinish() {
	if a.verbose {
		logger.Warn.Println("Аналитика финиширует")
	}
	a.Ping("100_lab_finish", nil)
}

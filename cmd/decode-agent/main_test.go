package main

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"gotest.tools/assert"
)

var mockStore []interface{}

type mockClient struct{}

func (mock *mockClient) IsConnected() bool {
	return true
}

func (mock *mockClient) IsConnectionOpen() bool {
	return true
}

func (mock *mockClient) Connect() MQTT.Token {
	return mockToken{}
}

func (mock *mockClient) Disconnect(quiesce uint) {}

func (mock *mockClient) Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token {
	if mockStore == nil {
		mockStore = make([]interface{}, 100)
	}
	mockStore = append(mockStore, payload)
	return mockToken{}
}

func (mock *mockClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	return mockToken{}
}

func (mock *mockClient) SubscribeMultiple(filters map[string]byte, callback MQTT.MessageHandler) MQTT.Token {
	return mockToken{}
}
func (mock *mockClient) Unsubscribe(topics ...string) MQTT.Token {
	return mockToken{}
}

func (mock *mockClient) AddRoute(topic string, callback MQTT.MessageHandler) {}

func (mock *mockClient) OptionsReader() MQTT.ClientOptionsReader {
	return MQTT.ClientOptionsReader{}
}

type mockToken struct{}

func (m mockToken) Wait() bool {
	return true
}

func (m mockToken) WaitTimeout(time.Duration) bool {
	return false
}

func (m mockToken) Error() error {
	return nil
}

func TestPlayback(t *testing.T) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// stub out mqtt method
	connectToMqtt = func(server string, clientid string, username string, password string) (MQTT.Client, error) {
		return &mockClient{}, nil
	}
	// read configuration
	cfg, err := cfgparser.Parse("./test/resources/config/playback.yaml")
	assert.NilError(t, err)
	playback(cfg)
	lastMsg, ok := mockStore[len(mockStore)-1].([]byte)
	assert.Assert(t, ok)
	lastStr := hex.EncodeToString(lastMsg)
	assert.Equal(t, strings.ToUpper(lastStr), "80142E4140049855C407A76D84C11CB2FD1488017FFFFFFFF00002EFFD7A37C14E8005800011823100082000103400480003035B7D5233D38000")
}

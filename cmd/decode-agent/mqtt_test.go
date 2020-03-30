package main

import (
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// mocks objects
type mockClient struct {
	PubMux    sync.RWMutex
	CbMux     sync.Mutex
	MockStore []interface{}
	CallBack  MQTT.MessageHandler
}

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
	mock.PubMux.Lock()
	mock.MockStore = append(mock.MockStore, payload)
	mock.PubMux.Unlock()
	return mockToken{}
}

func (mock *mockClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	mock.CbMux.Lock()
	mock.CallBack = callback
	mock.CbMux.Unlock()
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

type message struct {
	payload []byte
}

func (m *message) Duplicate() bool {
	return false
}

func (m *message) Qos() byte {
	return byte(0)
}

func (m *message) Retained() bool {
	return false
}

func (m *message) Topic() string {
	return "FOO"
}

func (m *message) MessageID() uint16 {
	return uint16(21321)
}

func (m *message) Payload() []byte {
	return m.payload
}

func (m *message) Ack() {}

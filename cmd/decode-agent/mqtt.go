package main

import (
	"crypto/tls"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/yh742/j2735-decoder/internal/cfgparser"
)

func generateClientID() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "host-" + string(rand.Int63()) + string(rand.Int63())
	}
	return hostname + "-" + strconv.Itoa(time.Now().Second())
}

func createMQTTClient(setting cfgparser.MqttSettings, callback MQTT.MessageHandler) (MQTT.Client, error) {
	auth := parseAuthFiles(setting.MQTTAuth)
	cli, err := connectToMqtt(setting.Server, setting.Clientid, auth, callback)
	if err != nil {
		log.Error().
			Err(err).
			Msg("cannot connect to mqtt server for publishing")
		return nil, err
	}
	log.Debug().Msgf("Connected to %s", setting.Server)
	return cli, nil
}

// Set this way so we can test this stub this out for test
var connectToMqtt = func(server string, clientid string, auth basicAuth, callback MQTT.MessageHandler) (MQTT.Client, error) {
	connOpts := MQTT.NewClientOptions()
	if strings.TrimSpace(server) != "" {
		connOpts.AddBroker(server)
	} else {
		return nil, errors.New("server string cannot be empty")
	}
	if auth.username != "" {
		log.Debug().Msg("username specified")
		connOpts.SetUsername(auth.username)
	}
	if auth.password != "" {
		log.Debug().Msg("password specified")
		connOpts.SetPassword(auth.password)
	}
	if strings.TrimSpace(clientid) == "" {
		clientid = generateClientID()
	}
	connOpts.SetClientID(clientid)
	connOpts.SetCleanSession(true)
	connOpts.DefaultPublishHandler = callback

	// skip tls config for now
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	// create client
	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

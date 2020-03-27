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
)

func generateClientID() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "host-" + string(rand.Int63())
	}
	return hostname + "-" + strconv.Itoa(time.Now().Second())
}

// Set this way so we can test this stub this out for test
var connectToMqtt = func(server string, clientid string, username string, password string) (MQTT.Client, error) {
	connOpts := MQTT.NewClientOptions().AddBroker(server).SetClientID(clientid).SetCleanSession(true)
	if strings.TrimSpace(server) != "" {
		connOpts.AddBroker(server)
	} else {
		return nil, errors.New("server string cannot be empty")
	}
	if username != "" {
		log.Debug().Msg("username specified")
		connOpts.SetUsername(username)
	}
	if password != "" {
		log.Debug().Msg("password specified")
		connOpts.SetPassword(password)
	}
	connOpts.SetClientID(clientid)
	connOpts.SetCleanSession(true)

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

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"github.com/yh742/j2735-decoder/pkg/decoder"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

var sig = make(chan os.Signal)
var testReady = make(chan bool, 1)

func messageHandler(msg MQTT.Message, pubClient MQTT.Client, pubCfg cfgparser.MqttSettings, format decoder.StringFormatType) {
	// log.Debug().Str("source", msg.Topic()).Msgf("%X", msg.Payload())
	var token MQTT.Token
	if format != decoder.PASS {
		decodedMsg, err := decoder.DecodeBytes(msg.Payload(), uint(len(msg.Payload())), format, msg.Topic())
		if err != nil {
			log.Error().Err(err).Msg("cannot decode msg format")
			return
		}
		token = pubClient.Publish(pubCfg.Topic, byte(pubCfg.Qos), false, decodedMsg)
	} else {
		token = pubClient.Publish(pubCfg.Topic, byte(pubCfg.Qos), false, msg.Payload())
	}
	token.Wait()
}

func bridge(cfg cfgparser.Config) {
	// print out bridge related information
	log.Info().
		Uint8("op.playbackcfg.format", uint8(cfg.Op.Format)).
		Send()

	// when user control-c out of the program
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// grab a connection to publish to
	pubClient, err := connectToMqtt(cfg.Publish.Server, cfg.Publish.Clientid, cfg.Publish.Username, cfg.Publish.Password)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot connect to mqtt server for publishing")
	}
	defer pubClient.Disconnect(500)
	log.Debug().Msgf("Connected to %s", cfg.Publish.Server)

	// grab a connection subscribe to
	subClient, err := connectToMqtt(cfg.Subscribe.Server, cfg.Subscribe.Clientid, cfg.Subscribe.Username, cfg.Subscribe.Password)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot connect to mqtt server for subscribing")
	}
	defer subClient.Disconnect(500)
	log.Debug().Msgf("Connected to %s", cfg.Subscribe.Server)
	subClient.Subscribe(cfg.Subscribe.Topic, cfg.Subscribe.Qos, func(client MQTT.Client, msg MQTT.Message) {
		messageHandler(msg, pubClient, cfg.Publish.MqttSettings, cfg.Op.Format)
	})
	testReady <- true
	log.Debug().Msg("waiting for ctr-c...")
	<-sig
	log.Debug().Msg("disconnecting...")
}

package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"github.com/yh742/j2735-decoder/pkg/decoder"
)

type bridge struct {
	cfg       cfgparser.Config
	pubClient MQTT.Client
	subClient MQTT.Client
	callBack  MQTT.MessageHandler
}

func newBridgeConnection(cfg cfgparser.Config) (*bridge, error) {
	// print out bridge related information
	log.Info().
		Uint8("op.playbackcfg.format", uint8(cfg.Op.Format)).
		Send()

	// init variables
	bdConn := bridge{}
	bdConn.cfg = cfg

	// grab a connection to publish to
	var err error
	bdConn.pubClient, err = createMQTTClient(cfg.Publish.MqttSettings, nil)
	if err != nil {
		return nil, err
	}

	// define a default mqtt message handler here
	bdConn.subClient, err = createMQTTClient(cfg.Subscribe.MqttSettings, nil)
	if err != nil {
		return nil, err
	}
	return &bdConn, nil
}

func (agt *bridge) startListening(callback MQTT.MessageHandler) bool {
	if callback != nil {
		agt.callBack = callback
	}
	token := agt.subClient.Subscribe(agt.cfg.Subscribe.Topic, agt.cfg.Subscribe.Qos, callback)
	return token.Wait()
}

func (agt *bridge) updateSettings(newSettings ExposedSettings) {
	if newSettings.PubTopic != agt.cfg.Publish.Topic && newSettings.PubTopic != "" {
		agt.switchTopic(agt.pubClient, agt.cfg.Publish.MqttSettings, newSettings.PubTopic)
		agt.cfg.Publish.Topic = newSettings.PubTopic
	}
	if newSettings.SubTopic != agt.cfg.Subscribe.Topic && newSettings.SubTopic != "" {
		agt.switchTopic(agt.subClient, agt.cfg.Subscribe.MqttSettings, newSettings.SubTopic)
		agt.cfg.Subscribe.Topic = newSettings.SubTopic
	}
	if newSettings.Format != agt.cfg.Op.Format && newSettings.Format != decoder.NA {
		agt.cfg.Op.Format = newSettings.Format
	}
}

func (agt *bridge) switchTopic(client MQTT.Client, currentSetting cfgparser.MqttSettings, newTopic string) {
	token := client.Unsubscribe(currentSetting.Topic).Wait()
	if !token {
		log.Error().Msgf("could not unsubscribe from topic %s", currentSetting.Topic)
	}
	token = client.Subscribe(newTopic, currentSetting.Qos, agt.callBack).Wait()
	if !token {
		log.Error().Msgf("could not subscribe to new topic %s", newTopic)
	}
}

func (agt *bridge) disconnect() {
	agt.pubClient.Disconnect(500)
	agt.subClient.Disconnect(500)
}

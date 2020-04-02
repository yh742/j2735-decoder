package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/yh742/j2735-decoder/internal/cfgparser"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func newAgent(mode cfgparser.CfgMode) agentRunner {
	switch mode {
	case cfgparser.Stream:
		return &streamAgent{}
	case cfgparser.Batch:
	}
	return nil
}

type agentRunner interface {
	run(cfgparser.Config, bool) error
	kill()
}

type streamAgent struct {
	killSig chan os.Signal
	bridge  *bridge
	block   bool
}

func (agt *streamAgent) run(cfg cfgparser.Config, block bool) error {
	agt.killSig = make(chan os.Signal)
	signal.Notify(agt.killSig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	bdConn, err := newBridgeConnection(cfg)
	if err != nil {
		return err
	}
	bdConn.startHTTPServer(8080)
	bdConn.startListening(func(client MQTT.Client, msg MQTT.Message) {
		mqttMessageHandler(bdConn.pubClient, msg, bdConn.cfg.Publish.MqttSettings, bdConn.cfg.Op.Format)
	})
	agt.bridge = bdConn
	agt.block = block
	if agt.block {
		<-agt.killSig
	}
	return nil
}

func (agt *streamAgent) kill() {
	if agt.block {
		agt.killSig <- syscall.SIGINT
	}
	agt.bridge.disconnect()
}

type batchAgent struct{}

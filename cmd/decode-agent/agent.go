package main

import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"github.com/yh742/j2735-decoder/pkg/decoder"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func newAgent(mode cfgparser.CfgMode) agentRunner {
	switch mode {
	case cfgparser.Stream:
		return &streamAgent{}
	case cfgparser.Batch:
		return &batchAgent{}
	}
	return nil
}

type agentRunner interface {
	run(cfgparser.Config, bool) error
	getBridge() *bridge
	kill()
}

type streamAgent struct {
	killSig chan os.Signal
	bridge  *bridge
	block   bool
}

func (agt *streamAgent) getBridge() *bridge {
	return agt.bridge
}

func (agt *streamAgent) run(cfg cfgparser.Config, block bool) error {
	agt.killSig = make(chan os.Signal)
	signal.Notify(agt.killSig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	bdConn, err := newBridgeConnection(cfg)
	if err != nil {
		return err
	}
	bdConn.startListening(func(client MQTT.Client, msg MQTT.Message) {
		log.Debug().Msgf("received msg: %v", msg)
		mqttMessageHandler(bdConn.pubClient, msg, bdConn.cfg.Publish.MqttSettings, bdConn.cfg.Op)
	})
	agt.bridge = bdConn
	agt.block = block
	if agt.block {
		log.Debug().Msg("waiting for stop signal")
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

type batchAgent struct {
	killSig chan os.Signal
	bridge  *bridge
	block   bool
	pMap    map[uint64]string
	setChan chan string
	msgIdx  uint64
}

func (agt *batchAgent) getBridge() *bridge {
	return agt.bridge
}

func (agt *batchAgent) run(cfg cfgparser.Config, block bool) error {
	agt.killSig = make(chan os.Signal)
	signal.Notify(agt.killSig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	agt.msgIdx = 0
	agt.setChan = make(chan string)
	agt.pMap = make(map[uint64]string)
	bdConn, err := newBridgeConnection(cfg)
	if err != nil {
		return err
	}
	bdConn.startListening(func(client MQTT.Client, msg MQTT.Message) {
		decodedMsg, err := decoder.DecodeBytes(msg.Payload(), bdConn.cfg.Op.Format, msg.Topic(), bdConn.cfg.Op.UseProtoBuf)
		if err != nil {
			log.Error().Err(err).Msg("cannot decode msg format")
			return
		}
		agt.setChan <- decodedMsg
	})
	agt.bridge = bdConn

	// frequency at which we publish the map
	ticker := time.Tick(time.Duration(cfg.Op.BatchCfg.PubFreq) * time.Millisecond)
	processMsg := func() error {
		for {
			select {
			case <-ticker:
				str := ""
				for _, v := range agt.pMap {
					str += v
				}
				if str != "" {
					bdConn.pubClient.Publish(bdConn.cfg.Publish.Topic, bdConn.cfg.Publish.Qos, false, str).Wait()
					agt.pMap = make(map[uint64]string)
				}
			case msg := <-agt.setChan:
				agt.pMap[agt.msgIdx] = msg
				atomic.AddUint64(&agt.msgIdx, 1)
			case <-agt.killSig:
				return nil
			}
		}
	}
	if block {
		return processMsg()
	}
	go processMsg()
	return nil
}

func (agt *batchAgent) kill() {
	if agt.block {
		agt.killSig <- syscall.SIGINT
	}
	agt.bridge.disconnect()
}

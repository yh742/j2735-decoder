package main

import (
	"context"
	"net/http"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"github.com/yh742/j2735-decoder/pkg/decoder"
)

// ExposedSettings are settings exposed through http
type ExposedSettings struct {
	PubTopic string
	SubTopic string
	Format   decoder.StringFormatType
}

type bridge struct {
	httpHupSig chan bool
	httpServer http.Server
	cfg        cfgparser.Config
	pubClient  MQTT.Client
	subClient  MQTT.Client
	callBack   MQTT.MessageHandler
}

func newBridgeConnection(cfg cfgparser.Config) (*bridge, error) {
	// print out bridge related information
	log.Info().
		Uint8("op.playbackcfg.format", uint8(cfg.Op.Format)).
		Send()

	// init variables
	bdConn := bridge{}
	bdConn.cfg = cfg
	bdConn.httpHupSig = make(chan bool)

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

func (agt *bridge) startHTTPServer(port int) {
	// listen on settings endpoint, exposing only sub/pub topics and format
	const URL = "/settings"

	httpAuth := parseAuthFiles(agt.cfg.Op.StreamCfg.HTTPAuth)
	log.Debug().Msgf("Username: '%s' Password: '%s'", httpAuth.username, httpAuth.password)
	router := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// GET methods
	router.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		getSettingHandler(w, r, &agt.cfg, httpAuth)
	}).Methods("GET")

	// PUT methods
	router.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		putSettingsHandler(w, r, httpAuth, func(eSetting ExposedSettings) {
			if eSetting.PubTopic != agt.cfg.Publish.Topic && eSetting.PubTopic != "" {
				agt.switchTopics(agt.pubClient, agt.cfg.Publish.MqttSettings, eSetting.PubTopic)
				agt.cfg.Publish.Topic = eSetting.PubTopic
			}
			if eSetting.SubTopic != agt.cfg.Subscribe.Topic && eSetting.SubTopic != "" {
				agt.switchTopics(agt.subClient, agt.cfg.Subscribe.MqttSettings, eSetting.SubTopic)
				agt.cfg.Subscribe.Topic = eSetting.SubTopic
			}
			if eSetting.Format != agt.cfg.Op.Format && eSetting.Format != decoder.NA {
				agt.cfg.Op.Format = eSetting.Format
			}
		})
	}).Methods("PUT")

	agt.httpServer = http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: handlers.CORS(originsOk, headersOk, methodsOk)(router),
	}
	go func() {
		if err := agt.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed at listenandserver()")
		}
		<-agt.httpHupSig
	}()
}

func (agt *bridge) switchTopics(client MQTT.Client, oldSetting cfgparser.MqttSettings, newTopic string) {
	token := client.Unsubscribe(oldSetting.Topic).Wait()
	if !token {
		log.Error().Msgf("could not unsubscribe from topic %s", oldSetting.Topic)
	}
	token = client.Subscribe(newTopic, oldSetting.Qos, agt.callBack).Wait()
	if !token {
		log.Error().Msgf("could not subscribe to new topic %s", newTopic)
	}
}

func (agt *bridge) disconnect() {
	agt.pubClient.Disconnect(500)
	agt.subClient.Disconnect(500)
	if agt.httpServer.Addr != "" {
		if err := agt.httpServer.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("can't shutdown http server properly")
		}
		log.Debug().Msg("http server teardown...")
		agt.httpHupSig <- true
		log.Debug().Msg("http server teardown finished...")
	}
}

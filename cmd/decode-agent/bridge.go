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
)

const pubURL = "/publish/setting"
const subURL = "/subscribe/setting"

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
	httpAuth := parseAuthFiles(agt.cfg.Op.StreamCfg.HTTPAuth)
	log.Debug().Msgf("Username: '%s' Password: '%s'", httpAuth.username, httpAuth.password)
	router := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// GET methods
	router.HandleFunc(pubURL, func(w http.ResponseWriter, r *http.Request) {
		getSettingHandler(w, r, &agt.cfg.Publish.MqttSettings, httpAuth)
	}).Methods("GET")
	router.HandleFunc(subURL, func(w http.ResponseWriter, r *http.Request) {
		getSettingHandler(w, r, &agt.cfg.Subscribe.MqttSettings, httpAuth)
	}).Methods("GET")

	// PUT methods
	router.HandleFunc(pubURL, func(w http.ResponseWriter, r *http.Request) {
		putSettingsHandler(w, r, agt.cfg.Publish.Topic, httpAuth, func(newTopic string) {
			agt.switchTopicCb(newTopic, pubURL)
		})
	}).Methods("PUT")
	router.HandleFunc(subURL, func(w http.ResponseWriter, r *http.Request) {
		putSettingsHandler(w, r, agt.cfg.Subscribe.Topic, httpAuth, func(newTopic string) {
			agt.switchTopicCb(newTopic, subURL)
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

func (agt *bridge) switchTopicCb(newTopic string, url string) {
	var client MQTT.Client
	var settings *cfgparser.MqttSettings
	switch url {
	case pubURL:
		client = agt.pubClient
		settings = &agt.cfg.Publish.MqttSettings
	case subURL:
		client = agt.subClient
		settings = &agt.cfg.Subscribe.MqttSettings
	}
	token := client.Unsubscribe(settings.Topic).Wait()
	if !token {
		log.Error().Msgf("could not unsubscribe from topic %s", settings.Topic)
	}
	settings.Topic = newTopic
	token = client.Subscribe(settings.Topic, settings.Qos, agt.callBack).Wait()
	if !token {
		log.Error().Msgf("could not subscribe to new topic %s", settings.Topic)
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

package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"github.com/yh742/j2735-decoder/pkg/decoder"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

var sig = make(chan os.Signal)
var testReady = make(chan bool, 1)

func messageHandler(msg MQTT.Message, pubClient MQTT.Client, pubCfg cfgparser.MqttSettings, format decoder.StringFormatType) {
	log.Debug().Str("source", msg.Topic()).Msgf("%X", msg.Payload())
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

func getSettings(w http.ResponseWriter, cfg *cfgparser.MqttSettings) {
	js, err := json.Marshal(cfg)
	if err != nil {
		log.Error().Err(err).Msg("cannot marshal json on http pub GET")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func putSettings(w http.ResponseWriter, r *http.Request, client MQTT.Client, cfg *cfgparser.MqttSettings, format decoder.StringFormatType) {
	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("cannot read http PUT body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	settings := cfgparser.MqttSettings{}
	err = json.Unmarshal(rbody, &settings)
	if err != nil {
		log.Error().Err(err).Msg("cannot read unmarshal into json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newTopic := strings.TrimSpace(settings.Topic)
	log.Debug().Msgf("unmarshaled results: %v", settings)
	if newTopic == "" || newTopic == cfg.Topic {
		log.Debug().Msg("content not updated")
		w.WriteHeader(204)
		return
	}
	log.Debug().Msgf("old config: %v", cfg)
	token := client.Unsubscribe(cfg.Topic)
	token.Wait()
	cfg.Topic = newTopic
	token = client.Subscribe(cfg.Topic, cfg.Qos, nil)
	token.Wait()
	log.Debug().Msgf("updated config: %v", cfg)
	w.WriteHeader(200)
}

func httpHandler(w http.ResponseWriter, r *http.Request, client MQTT.Client,
	cfg *cfgparser.MqttSettings, auth basicAuth, format decoder.StringFormatType) {
	if !checkBasicHTTPAuth(r, auth) {
		http.Error(w, "unable to verify identity", http.StatusForbidden)
		return
	}
	switch r.Method {
	case "GET":
		getSettings(w, cfg)
	case "PUT":
		putSettings(w, r, client, cfg, format)
	}
}

func bridge(cfg cfgparser.Config) {
	// print out bridge related information
	log.Info().
		Uint8("op.playbackcfg.format", uint8(cfg.Op.Format)).
		Send()

	// when user control-c out of the program
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// grab a connection to publish to
	pubAuth := parseAuthFiles(cfg.Publish.MQTTAuth)
	pubClient, err := connectToMqtt(cfg.Publish.Server, cfg.Publish.Clientid, pubAuth, nil)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot connect to mqtt server for publishing")
	}
	defer pubClient.Disconnect(500)
	log.Debug().Msgf("Connected to %s", cfg.Publish.Server)

	// grab a connection subscribe to
	subAuth := parseAuthFiles(cfg.Publish.MQTTAuth)
	subClient, err := connectToMqtt(cfg.Subscribe.Server, cfg.Subscribe.Clientid, subAuth,
		func(client MQTT.Client, msg MQTT.Message) {
			messageHandler(msg, pubClient, cfg.Publish.MqttSettings, cfg.Op.Format)
		})
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot connect to mqtt server for subscribing")
	}
	defer subClient.Disconnect(500)
	log.Debug().Msgf("Connected to %s", cfg.Subscribe.Server)
	// create a client subscription
	subClient.Subscribe(cfg.Subscribe.Topic, cfg.Subscribe.Qos, nil)

	// create a http handler
	httpAuth := parseAuthFiles(cfg.Op.BridgeCfg.HTTPAuth)
	log.Debug().Msgf("Username: '%s' Password: '%s'", httpAuth.username, httpAuth.password)
	router := mux.NewRouter()

	// setup headers and allow CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	router.HandleFunc("/publish/setting", func(w http.ResponseWriter, r *http.Request) {
		httpHandler(w, r, pubClient, &cfg.Publish.MqttSettings, httpAuth, cfg.Op.Format)
	}).Methods("GET", "PUT")
	router.HandleFunc("/subscribe/setting", func(w http.ResponseWriter, r *http.Request) {
		httpHandler(w, r, subClient, &cfg.Subscribe.MqttSettings, httpAuth, cfg.Op.Format)
	}).Methods("GET", "PUT")

	// make sure we can gracefully exit http server
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handlers.CORS(originsOk, headersOk, methodsOk)(router),
	}
	go func(wg *sync.WaitGroup) {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed at listenandserver()")
		}
		wg.Done()
	}(httpServerExitDone)

	// notify test it is ready to run
	testReady <- true
	log.Debug().Msg("waiting for ctr-c...")
	<-sig

	// shutdown http server
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("can't shutdown http server properly")
	}
	httpServerExitDone.Wait()
	log.Debug().Msg("disconnecting...")
}

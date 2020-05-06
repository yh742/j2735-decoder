package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"github.com/yh742/j2735-decoder/pkg/decoder"
)

func mqttMessageHandler(pubClient MQTT.Client, msg MQTT.Message, pubCfg cfgparser.MqttSettings, format decoder.StringFormatType) {
	log.Debug().Str("source", msg.Topic()).Msgf("%X", msg.Payload())
	var token MQTT.Token
	if format != decoder.PASS {
		decodedMsg, err := decoder.DecodeBytes(msg.Payload(), uint(len(msg.Payload())), format, msg.Topic())
		if err != nil {
			log.Error().Err(err).Msg("cannot decode msg format")
			return
		}
		token = pubClient.Publish(pubCfg.Topic, byte(pubCfg.Qos), false, decodedMsg)
		log.Trace().Msgf("decoded message in %d: %s", format, decodedMsg)
	} else {
		token = pubClient.Publish(pubCfg.Topic, byte(pubCfg.Qos), false, msg.Payload())
		log.Trace().Msgf("decoded message in %d: %s", format, msg.Payload())
	}
	token.Wait()
}

func getSettingHandler(w http.ResponseWriter, r *http.Request, cfg cfgparser.Config, auth basicAuth) {
	if !checkBasicHTTPAuth(r, auth) {
		http.Error(w, "unable to verify identity", http.StatusForbidden)
		return
	}
	// only expose mutable variables
	js, err := json.Marshal(ExposedSettings{
		PubTopic: cfg.Publish.Topic,
		SubTopic: cfg.Subscribe.Topic,
		Format:   cfg.Op.Format,
	})
	if err != nil {
		log.Error().Err(err).Msg("cannot marshal json on http pub GET")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func putSettingsHandler(w http.ResponseWriter, r *http.Request, auth basicAuth, updateCallback func(ExposedSettings)) {
	if !checkBasicHTTPAuth(r, auth) {
		http.Error(w, "unable to verify identity", http.StatusForbidden)
		return
	}
	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("cannot read http PUT body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	settings := ExposedSettings{}
	err = json.Unmarshal(rbody, &settings)
	if err != nil {
		log.Error().Err(err).Msg("cannot read unmarshal into json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug().Msgf("unmarshaled results: %v", settings)
	updateCallback(settings)
	w.WriteHeader(200)
}

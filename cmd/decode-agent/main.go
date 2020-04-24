package main

import (
	"flag"

	"github.com/yh742/j2735-decoder/internal/cfgparser"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// find the config file path
	cfgPath := flag.String("cfg", "/etc/decode-agent/config.yaml", "Config path for decode-agent")
	logLevel := flag.Int("loglevel", 1, "Set log level (trace=-1, debug=0, info=1, warn=2, error=3)")
	flag.Parse()

	// set log time format
	log.Logger = log.With().Caller().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.Level(*logLevel))

	// parse the config file
	cfg, err := cfgparser.Parse(*cfgPath)
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("Error occured reading config file at %s", *cfgPath)
	}

	// print out configuration
	log.Info().
		Str("publish.server", cfg.Publish.Server).
		Str("publish.topic", cfg.Publish.Topic).
		Uint8("publish.qos", cfg.Publish.Qos).
		Str("publish.clientid", cfg.Publish.Clientid).
		Str("publish.mqttauth", cfg.Publish.MQTTAuth).
		Str("subscribe.server", cfg.Subscribe.Server).
		Str("subscribe.topic", cfg.Subscribe.Topic).
		Uint8("subscribe.qos", cfg.Subscribe.Qos).
		Str("subscribe.clientid", cfg.Subscribe.Clientid).
		Str("subscribe.mqttauth", cfg.Subscribe.MQTTAuth).
		Uint8("op.mode", uint8(cfg.Op.Mode)).
		Send()

	// launch decode-agent as different modes
	// supported modes: playback, bridge
	// TODO: batching messages
	var agt agentRunner
	switch cfg.Op.Mode {
	case cfgparser.Playback:
		playback(cfg)
	default:
		agt = newAgent(cfg.Op.Mode)
	}
	if agt != nil {
		agt.run(cfg, true)
	}
	log.Fatal().Msg("agent type not supported")
}
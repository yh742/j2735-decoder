package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/yh742/j2735-decoder/internal/cfgparser"

	"github.com/rs/zerolog/log"
)

type playbackHandler func(data []byte)

func playLogFile(file *os.File, handler playbackHandler, loop bool, pubFreq uint) {
	reader := bufio.NewReader(file)
	lineCnt := 0
	for true {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Error().
				Err(err).
				Msgf("something occured while reading playback file line #%d", lineCnt)
			continue
		}
		if err == io.EOF {
			log.Debug().Msg("reached the end of file ...")
			if loop {
				file.Seek(0, 0)
				lineCnt = 0
				continue
			}
			break
		}
		var objMap map[string]interface{}
		if err = json.Unmarshal([]byte(line), &objMap); err != nil {
			log.Fatal().
				Err(err).
				Msgf("Unable to unmarshal json at line %d", lineCnt)
		}
		str, ok := objMap["message"].(string)
		if !ok {
			log.Fatal().
				Err(errors.New("Conversion Error")).
				Msgf("JSON message field is not populated with the right type at line %d", lineCnt)
		}
		log.Debug().Msgf("%s", str)
		data, err := hex.DecodeString(str)
		if err != nil {
			log.Warn().
				Err(err).
				Msgf("Could not decode hexstring at line %d, %s", lineCnt, str)
		}
		handler(data)
		lineCnt++
		time.Sleep(time.Duration(pubFreq) * time.Millisecond)
	}
}

func playback(cfg cfgparser.Config) {
	// print out playback configuration
	log.Info().
		Str("op.playbackcfg.file", cfg.Op.PlaybackCfg.File).
		Bool("op.playbackcfg.loop", cfg.Op.PlaybackCfg.Loop).
		Uint8("op.playbackcfg.format", uint8(cfg.Op.Format)).
		Send()

	// read in playback file
	file, err := os.Open(cfg.Op.PlaybackCfg.File)
	defer file.Close()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error reading playback file!")
	}

	// connect to MQTT service
	if cfg.Publish.Clientid == "" {
		// clientid is needed and needs to be unique
		cfg.Publish.Clientid = generateClientID()
	}
	auth := parseAuthFiles(cfg.Publish.MQTTAuth)
	pubClient, err := connectToMqtt(cfg.Publish.Server, cfg.Publish.Clientid, auth, nil)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot connect to mqtt server")
	}

	// read from file
	playLogFile(file, func(data []byte) {
		pubClient.Publish(cfg.Publish.Topic, cfg.Publish.Qos, false, data)
	}, cfg.Op.PlaybackCfg.Loop, cfg.Op.PlaybackCfg.PubFreq)
}

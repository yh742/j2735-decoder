package cfgparser

import (
	"strings"

	"github.com/yh742/j2735-decoder/pkg/decoder"
)

// Uint8Enum is an interface uint8 enums use to parse.
// Interfaces should be defined in packages that use it.
type Uint8Enum interface {
	ParseString(string) (uint8, bool)
}

// CfgMode represents decoder modes
type CfgMode uint8

const (
	// Batch is used for batch publishing decoded msgs
	Batch CfgMode = iota
	// Stream is used for publishing decoded msgs
	Stream
	// Playback is used for replaying uper msgs from a pre-recorded file
	Playback
	// Passthrough is used for bridging msgs (no decoding)
	Passthrough
)

// UnmarshalYAML for CfgMode type
func (mode *CfgMode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var modeStr string
	if err := unmarshal(&modeStr); err != nil {
		return err
	}
	switch strings.ToLower(modeStr) {
	case "batch":
		*mode = Batch
	case "playback":
		*mode = Playback
	case "passthrough":
		*mode = Passthrough
	case "stream":
		*mode = Stream
	}
	return nil
}

// ParseString returns number based on string representation
func (mode *CfgMode) ParseString(str string) (uint8, bool) {
	m := map[string]uint8{"batch": 0, "stream": 1, "playback": 2, "passthrough": 3}
	val, ok := m[str]
	return val, ok
}

// Config represents the yaml file fields
type Config struct {
	Subscribe struct {
		MqttSettings `yaml:",inline"` // should be uniform with publish in terms of naming convention for env vars
	}
	Publish struct {
		MqttSettings `yaml:",inline"`
	}
	Op struct {
		Mode        CfgMode
		BatchCfg    BatchConfig    `yaml:"batchconfig"`
		PlaybackCfg PlaybackConfig `yaml:"playbackconfig"`
	}
}

// PlaybackConfig are settings used for playback mode
type PlaybackConfig struct {
	File   string
	Loop   bool
	Format decoder.StringFormatType
}

// BatchConfig are settings used for batch mode
type BatchConfig struct {
	Pubfreq uint
	Expiry  uint
	Format  decoder.StringFormatType
}

// MqttSettings are settings used for batch mode
type MqttSettings struct {
	Clientid string
	Username string
	Password string
	Server   string
	Topic    string
	Qos      uint8
}

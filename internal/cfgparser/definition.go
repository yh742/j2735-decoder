package cfgparser

// Config represents the yaml file fields
type Config struct {
	Hostname  string
	Subscribe struct {
		MqttSettings `yaml:",inline"` // should be uniform with publish in terms of naming convention for env vars
	}
	Publish struct {
		MqttSettings `yaml:",inline"`
		Format       string
	}
	Auth struct {
		Clientid string
		Username string
		Password string
	}
	Op struct {
		Mode     string
		Settings struct {
			BatchConfig    `yaml:",inline"`
			PlaybackConfig `yaml:",inline"`
		}
	}
}

// PlaybackConfig is used for playback mode
type PlaybackConfig struct {
	File string
}

// BatchConfig is used for batch mode
type BatchConfig struct {
	Pubfreq uint
	Expiry  uint
}

type MqttSettings struct {
	Server string
	Topic  string
	Qos    uint8
}

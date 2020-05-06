package cfgparser_test

import (
	"os"
	"path"
	"testing"

	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"github.com/yh742/j2735-decoder/pkg/decoder"
	"gotest.tools/assert"
)

func TestParseYaml(t *testing.T) {
	testCfg := "config.yaml"
	t.Logf("Parsing '%s'... ", testCfg)
	cfgs, err := cfgparser.Parse(path.Join("./test/resources/", testCfg))
	assert.NilError(t, err)
	for _, cfg := range cfgs {
		assert.Equal(t, cfg.Subscribe.Qos, uint8(3))
		assert.Equal(t, cfg.Publish.Topic, "vz/outputs")
		assert.Equal(t, cfg.Publish.Clientid, "111")
		switch cfg.Op.Mode {
		case cfgparser.Stream:
			assert.Equal(t, cfg.Name, "streamConfig")
			assert.Equal(t, cfg.Op.HTTPAuth, "/etc/decode-agent/http-psw")
			assert.Equal(t, cfg.Op.Format, decoder.JSON)
		case cfgparser.Playback:
			assert.Equal(t, cfg.Name, "playbackConfig")
			assert.Equal(t, cfg.Op.PlaybackCfg.File, "./playback.txt")
			assert.Equal(t, cfg.Op.PlaybackCfg.Loop, true)
			assert.Equal(t, cfg.Op.Format, decoder.JSON)
		default:
			t.Log("No suitable operation mode found!")
			t.Fail()
		}
	}
}

// // TestEnvVariables tests if environment variables overrides settings
func TestEnvVariables(t *testing.T) {
	os.Setenv("PLAYBACKCONFIG_PUBLISH_MQTTSETTINGS_CLIENTID", "222")
	os.Setenv("PLAYBACKCONFIG_SUBSCRIBE_MQTTSETTINGS_SERVER", "subscribe.net")
	os.Setenv("PLAYBACKCONFIG_OP_MODE", "PLAYBACK")
	os.Setenv("PLAYBACKCONFIG_OP_FORMAT", "XML")
	os.Setenv("PLAYBACKCONFIG_OP_HTTPAUTH", "/etc/blah")
	cfgs, err := cfgparser.Parse(path.Join("./test/resources/config.yaml"))
	assert.NilError(t, err)

	// these should match the environment variables you set above
	assert.Equal(t, cfgs[0].Publish.Clientid, "222")
	assert.Equal(t, cfgs[0].Subscribe.Server, "subscribe.net")
	assert.Equal(t, cfgs[0].Op.Mode, cfgparser.Playback)
	assert.Equal(t, cfgs[0].Op.HTTPAuth, "/etc/blah")
	assert.Equal(t, cfgs[0].Op.Format, decoder.XML)
}

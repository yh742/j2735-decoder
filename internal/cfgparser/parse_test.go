package cfgparser_test

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"gotest.tools/assert"
)

func TestParseYaml(t *testing.T) {
	testArray := [...]string{"config.yaml", "config2.yaml"}
	for _, item := range testArray {
		t.Logf("Parsing '%s'... ", item)
		cfg, err := cfgparser.Parse(path.Join("./test/resources/", item))
		assert.NilError(t, err)
		assert.Equal(t, cfg.Hostname, "xyz.com")
		assert.Equal(t, cfg.Subscribe.Qos, uint8(3))
		assert.Equal(t, cfg.Publish.Topic, "vz/outputs")
		assert.Equal(t, cfg.Publish.Format, "json")
		assert.Equal(t, cfg.Auth.Clientid, "111")
		switch cfg.Op.Mode {
		case "batch":
			assert.Equal(t, cfg.Op.Settings.Pubfreq, uint(200))
		case "playback":
			assert.Equal(t, cfg.Op.Settings.File, "./playback.txt")
		default:
			log.Fatal("No suitable operation mode found!")
		}
	}
}

// TestEnvVariables tests if environment variables overrides settings
func TestEnvVariables(t *testing.T) {
	os.Setenv("HOSTNAME", "abc.com")
	os.Setenv("SUBSCRIBE_MQTTSETTINGS_SERVER", "subscribe.net")
	os.Setenv("PUBLISH_FORMAT", "xml")
	os.Setenv("AUTH_CLIENTID", "222")
	os.Setenv("OP_MODE", "PLAYBACK")
	os.Setenv("OP_SETTINGS_BATCHCONFIG_PUBFREQ", "900")
	cfg, err := cfgparser.Parse(path.Join("./test/resources/config.yaml"))
	assert.NilError(t, err)
	t.Logf("%v", cfg)
	assert.Equal(t, cfg.Hostname, "abc.com")
	assert.Equal(t, cfg.Subscribe.Server, "subscribe.net")
	assert.Equal(t, cfg.Publish.Format, "xml")
	assert.Equal(t, cfg.Auth.Clientid, "222")
	assert.Equal(t, cfg.Op.Settings.BatchConfig.Pubfreq, uint(900))
}

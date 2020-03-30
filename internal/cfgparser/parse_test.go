package cfgparser_test

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"github.com/yh742/j2735-decoder/pkg/decoder"
	"gotest.tools/assert"
)

func TestParseYaml(t *testing.T) {
	testArray := [...]string{"batch.yaml", "playback.yaml"}
	for _, item := range testArray {
		t.Logf("Parsing '%s'... ", item)
		cfg, err := cfgparser.Parse(path.Join("./test/resources/", item))
		assert.NilError(t, err)
		assert.Equal(t, cfg.Subscribe.Qos, uint8(3))
		assert.Equal(t, cfg.Publish.Topic, "vz/outputs")
		assert.Equal(t, cfg.Publish.Clientid, "111")
		switch cfg.Op.Mode {
		case cfgparser.Batch:
			assert.Equal(t, cfg.Op.BatchCfg.Pubfreq, uint(200))
			assert.Equal(t, cfg.Op.Format, decoder.JSON)
		case cfgparser.Playback:
			assert.Equal(t, cfg.Op.PlaybackCfg.File, "./playback.txt")
			assert.Equal(t, cfg.Op.PlaybackCfg.Loop, true)
			assert.Equal(t, cfg.Op.Format, decoder.JSON)
		default:
			log.Fatal("No suitable operation mode found!")
		}
	}
}

// TestEnvVariables tests if environment variables overrides settings
func TestEnvVariables(t *testing.T) {
	os.Setenv("PUBLISH_MQTTSETTINGS_CLIENTID", "222")
	os.Setenv("SUBSCRIBE_MQTTSETTINGS_SERVER", "subscribe.net")
	os.Setenv("OP_MODE", "PLAYBACK")
	os.Setenv("OP_FORMAT", "XML")
	os.Setenv("OP_BATCHCFG_PUBFREQ", "900")
	cfg, err := cfgparser.Parse(path.Join("./test/resources/batch.yaml"))
	assert.NilError(t, err)

	// these should match the environment variables you set above
	assert.Equal(t, cfg.Publish.Clientid, "222")
	assert.Equal(t, cfg.Subscribe.Server, "subscribe.net")
	assert.Equal(t, cfg.Op.Mode, cfgparser.Playback)
	assert.Equal(t, cfg.Op.BatchCfg.Pubfreq, uint(900))
	assert.Equal(t, cfg.Op.Format, decoder.XML)
}

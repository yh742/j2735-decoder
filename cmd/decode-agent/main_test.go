package main

import (
	"bufio"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"testing"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"gotest.tools/assert"
)

var mc *mockClient

func TestMain(m *testing.M) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// stub out mqtt method
	connectToMqtt = func(server string, clientid string, username string, password string) (MQTT.Client, error) {
		return mc, nil
	}
	os.Exit(m.Run())
}

func TestPlayback(t *testing.T) {
	// read configuration
	mc = &mockClient{}
	cfg, err := cfgparser.Parse("./test/resources/config/playback.yaml")
	assert.NilError(t, err)
	playback(cfg)
	mc.PubMux.RLock()
	lastMsg, ok := mc.MockStore[len(mc.MockStore)-1].([]byte)
	mc.PubMux.RUnlock()
	assert.Assert(t, ok)
	lastStr := hex.EncodeToString(lastMsg)
	assert.Equal(t,
		strings.ToUpper(lastStr),
		"80142E4140049855C407A76D84C11CB2FD1488017FFFFFFFF00002EFFD7A37C14E8005800011823100082000103400480003035B7D5233D38000")
}

func TestBridgePassthrough(t *testing.T) {
	// read configuration
	mc = &mockClient{}
	cfg, err := cfgparser.Parse("./test/resources/config/bridge-passthrough.yaml")
	assert.NilError(t, err)
	go func() {
		bridge(cfg)
	}()
	<-testReady
	file, err := os.Open("./test/resources/logs/bsm-sample.log")
	defer file.Close()
	assert.NilError(t, err)
	reader := bufio.NewReader(file)
	for true {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			t.Log("EOF reached")
			break
		}
		assert.NilError(t, err)
		mc.CallBack(mc, &message{
			payload: []byte(line),
		})
	}
	sig <- syscall.SIGINT
	mc.PubMux.RLock()
	lastMsg, ok := mc.MockStore[len(mc.MockStore)-1].([]byte)
	mc.PubMux.RUnlock()
	assert.Assert(t, ok)
	file.Seek(0, 0)
	data, err := ioutil.ReadAll(file)
	assert.NilError(t, err)
	assert.Equal(t, string(lastMsg[len(lastMsg)-20:]), string(data[len(data)-20:]))
}

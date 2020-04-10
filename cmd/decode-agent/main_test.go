package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yh742/j2735-decoder/internal/cfgparser"
	"github.com/yh742/j2735-decoder/pkg/decoder"
	"gotest.tools/assert"
)

var mc *mockClient

func TestMain(m *testing.M) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.With().Caller().Logger()
	// stub out mqtt method
	connectToMqtt = func(server string, clientid string, auth basicAuth, callback MQTT.MessageHandler) (MQTT.Client, error) {
		if callback != nil {
			mc.CbMux.Lock()
			mc.CallBack = callback
			mc.CbMux.Unlock()
		}
		return mc, nil
	}
	os.Exit(m.Run())
}

func TestPlayback(t *testing.T) {
	// read configuration
	mc = &mockClient{}
	const matchString = "80142E4140049855C407A76D84C11CB2FD1488017FFFFFFFF00002EFFD7A37C14E8005800011823100082000103400480003035B7D5233D38000"
	cfg, err := cfgparser.Parse("./test/resources/config/playback.yaml")
	assert.NilError(t, err)
	playback(cfg)
	mc.PubMux.RLock()
	lastMsg, ok := mc.MockStore[len(mc.MockStore)-1].([]byte)
	mc.PubMux.RUnlock()
	assert.Assert(t, ok)
	lastStr := hex.EncodeToString(lastMsg)
	assert.Equal(t, strings.ToUpper(lastStr), matchString)
}

func TestStreaming(t *testing.T) {
	// this test is noisy
	cfgs := map[string]string{
		"bridge-passthrough.yaml": "80142E4140049855C407A76D84C11CB2FD1488017FFFFFFFF00002EFFD7A37C14E8005800011823100082000103400480003035B7D5233D38000",
		"bridge-decode.yaml":      `{"MessageFrame": {"source": "FOO", "messageId": "20", "value": {"BasicSafetyMessage": {"partII": {"PartIIcontent": [{"partII-Id": "0", "partII-Value": {"VehicleSafetyExtensions": {"lights": "000000000"}}}, {"partII-Id": "1", "partII-Value": {"SpecialVehicleExtensions": {"vehicleAlerts": {"multi": {"unavailable": ""}, "sspRights": "0", "sirenUse": {"notInUse": ""}, "lightsUse": {"notInUse": ""}}}}}]}, "coreData": {"msgCnt": "5", "lat": "422977666", "speed": "0", "heading": "751", "size": {"length": "70", "width": "35"}, "id": "00126157", "secMark": "4126", "long": "-837015510", "elev": "2", "accuracy": {"semiMajor": "255", "semiMinor": "255", "orientation": "65535"}, "angle": "127", "accelSet": {"long": "-45", "lat": "-15", "vert": "-49", "yaw": "6"}, "brakes": {"wheelBrakes": "10000", "traction": {"unavailable": ""}, "abs": {"unavailable": ""}, "scs": {"unavailable": ""}, "brakeBoost": {"unavailable": ""}, "auxBrakes": {"unavailable": ""}}, "transmission": {"unavailable": ""}}}}}}`,
	}
	for key, value := range cfgs {
		mc = &mockClient{}
		cfg, err := cfgparser.Parse(path.Join("./test/resources/config/", key))
		assert.NilError(t, err)
		// launch bridge asynchronously
		sa := streamAgent{}
		sa.run(cfg, false)
		file, err := os.Open("./test/resources/logs/bsm-sample.log")
		defer file.Close()
		assert.NilError(t, err)
		playLogFile(file, func(data []byte) {
			mc.CallBack(mc, &message{
				payload: data,
			})
		}, false, cfg.Op.PlaybackCfg.PubFreq)
		sa.kill()
		if key == "bridge-passthrough.yaml" {
			mc.PubMux.RLock()
			lastMsg, ok := mc.MockStore[len(mc.MockStore)-1].([]byte)
			mc.PubMux.RUnlock()
			assert.Assert(t, ok)
			assert.Equal(t, hex.EncodeToString(lastMsg), strings.ToLower(value))
		} else {
			mc.PubMux.RLock()
			lastMsg, ok := mc.MockStore[len(mc.MockStore)-1].(string)
			mc.PubMux.RUnlock()
			assert.Assert(t, ok)
			ok, err := decoder.AreEqualJSON(lastMsg, value)
			assert.NilError(t, err)
			assert.Assert(t, ok)
		}
	}
}

func TestBatch(t *testing.T) {
	// this test is noisy
	mc = &mockClient{}
	cfg, err := cfgparser.Parse(path.Join("./test/resources/config/", "bridge-batch.yaml"))
	assert.NilError(t, err)
	// launch bridge asynchronously
	sa := batchAgent{}
	sa.run(cfg, false)
	file, err := os.Open("./test/resources/logs/bsm-sample.log")
	defer file.Close()
	assert.NilError(t, err)
	playLogFile(file, func(data []byte) {
		mc.CallBack(mc, &message{
			payload: data,
		})
	}, false, 0)
	sa.kill()
	mc.PubMux.RLock()
	_, ok := mc.MockStore[len(mc.MockStore)-1].(string)
	mc.PubMux.RUnlock()
	assert.Assert(t, ok)
}

func TestHttpGetPut(t *testing.T) {
	// read configuration
	mc = &mockClient{}
	cfg, err := cfgparser.Parse("./test/resources/config/bridge-passthrough.yaml")
	assert.NilError(t, err)
	// launch bridge asynchronously
	sa := streamAgent{}
	sa.run(cfg, false)
	client := &http.Client{}
	// check GET calls
	req, err := http.NewRequest("GET", "http://localhost:8080/settings", nil)
	assert.NilError(t, err)
	req.SetBasicAuth("admin", "admin")
	resp, err := client.Do(req)
	assert.NilError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	// check PUT calls
	reqBody, err := json.Marshal(map[string]string{
		"topic": "test/test",
	})
	assert.NilError(t, err)
	req, err = http.NewRequest("PUT", "http://localhost:8080/settings", bytes.NewBuffer(reqBody))
	assert.NilError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin")
	resp, err = client.Do(req)
	assert.NilError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	// check error condition
	req.Body = ioutil.NopCloser(strings.NewReader("blahblah"))
	req.ContentLength = int64(len("blahblah"))
	resp, err = client.Do(req)
	assert.NilError(t, err)
	assert.Equal(t, resp.StatusCode, 500)
	sa.kill()
}

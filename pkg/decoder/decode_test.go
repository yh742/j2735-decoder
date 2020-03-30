package decoder_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/yh742/j2735-decoder/pkg/decoder"
	"gotest.tools/assert"
)

const (
	InputFolder      = "./test/resource/in"
	XmlOutputFolder  = "./test/resource/out/xml"
	JsonOutputFolder = "./test/resource/out/json"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func TestUperToXmlStringConversion(t *testing.T) {
	// filenames
	testArray := [...]string{"bsm1", "bsm2", "bsm3", "psm1", "spat"}
	for _, item := range testArray {
		t.Logf("decoding '%s' to xml", item+".uper")
		// open input
		uperFile, err := os.Open(path.Join(InputFolder, item+".uper"))
		defer uperFile.Close()
		assert.NilError(t, err)
		// read expected output in xml
		xmlFile, err := os.Open(path.Join(XmlOutputFolder, item+".xml"))
		defer xmlFile.Close()
		uperBytes, err := ioutil.ReadAll(uperFile)
		assert.NilError(t, err)
		// decode
		decodedMsg, err := decoder.DecodeBytes(uperBytes, uint(len(uperBytes)), decoder.XML, "FOO")
		assert.NilError(t, err)
		xmlBytes, err := ioutil.ReadAll(xmlFile)
		assert.NilError(t, err)
		xmlString := fmt.Sprintf("%s", xmlBytes)
		assert.Equal(t, decodedMsg, xmlString)
	}
}

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}
	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}
	return reflect.DeepEqual(o1, o2), nil
}

func TestUperToJsonStringConversion(t *testing.T) {
	// filenames
	testArray := [...]string{"bsm1", "bsm2", "bsm3", "psm1", "spat"}
	for _, item := range testArray {
		t.Logf("decoding '%s' to json", item+".uper")
		// open input
		uperFile, err := os.Open(path.Join(InputFolder, item+".uper"))
		defer uperFile.Close()
		assert.NilError(t, err)
		// read expected output in json
		jsonFile, err := os.Open(path.Join(JsonOutputFolder, item+".json"))
		defer jsonFile.Close()
		uperBytes, err := ioutil.ReadAll(uperFile)
		assert.NilError(t, err)
		// decode
		decodedMsg, err := decoder.DecodeBytes(uperBytes, uint(len(uperBytes)), decoder.JSON, "FOO")
		assert.NilError(t, err)
		jsonBytes, err := ioutil.ReadAll(jsonFile)
		assert.NilError(t, err)
		jsonString := fmt.Sprintf("%s", jsonBytes)
		t.Logf("%s", decodedMsg)
		eq, err := AreEqualJSON(jsonString, decodedMsg)
		assert.NilError(t, err)
		assert.Assert(t, eq)
	}
}

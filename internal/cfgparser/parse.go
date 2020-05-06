package cfgparser

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// Parse takes in a yaml configuration file and unmarshals it.
// Configuration file can be overriden by environment variables.
// Returns a Config struct and error.
func Parse(configPath string) ([]Config, error) {
	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	bytesRead := bytes.NewReader(data)
	var cfgs []Config
	var cfg Config
	yamlDec := yaml.NewDecoder(bytesRead)
	for yamlDec.Decode(&cfg) == nil {
		overrideConfig(&cfg, cfg.Name)
		cfgs = append(cfgs, cfg)
	}
	return cfgs, nil
}

// overrideConfig recurses the config struct to look for corresponding env var.
// If env vars are set, it will replace the config value read from the yaml config file (e.g. Publish.Format => PUBLISH_FORMAT).
// Embedded structs will take on the type's name (e.g. Subscribe.Qos => SUBSCRIBE_MQTTSETTINGS_QOS).
func overrideConfig(obj interface{}, parentName string) {
	s := reflect.Indirect(reflect.ValueOf(obj))
	for i := 0; i < s.NumField(); i++ {
		t := s.Type().Field(i)
		v := s.Field(i)
		if fieldIsExported(t) {
			name := t.Name
			if parentName != "" {
				name = parentName + "_" + t.Name
			}
			if v.Kind() == reflect.Struct {
				overrideConfig(v.Addr().Interface(), name)
			} else {
				varName := strings.ToUpper(name)
				if os.Getenv(varName) != "" && v.CanSet() {
					// only 3 types exist in current struct, add more in needed
					switch v.Kind() {
					case reflect.String:
						v.SetString(os.Getenv(varName))
					case reflect.Bool:
						b, err := strconv.ParseBool(os.Getenv(varName))
						if err != nil {
							log.Warn().Msgf("Cannot convert Bool value from '%s'", os.Getenv(varName))
							break
						}
						v.SetBool(b)
					case reflect.Int8:
						num, err := strconv.ParseInt(os.Getenv(varName), 10, 8)
						if err != nil {
							log.Warn().Msgf("Cannot convert Int8 value from '%s'", os.Getenv(varName))
							break
						}
						v.SetInt(num)
					case reflect.Uint:
						num, err := strconv.ParseUint(os.Getenv(varName), 10, 64)
						if err != nil {
							log.Warn().Msgf("Cannot convert UInt value from '%s'", os.Getenv(varName))
							break
						}
						v.SetUint(num)
					case reflect.Uint8:
						num, err := strconv.ParseUint(os.Getenv(varName), 10, 8)
						if err != nil {
							val, ok := v.Addr().Interface().(Uint8Enum)
							if !ok {
								log.Warn().Msgf("Cannot get interface from '%s'", os.Getenv(varName))
								break
							}
							temp, ok := val.ParseString(strings.ToLower(os.Getenv(varName)))
							if !ok {
								log.Warn().Msgf("Cannot convert UInt8 value from '%s'", os.Getenv(varName))
								break
							}
							num = uint64(temp)
						}
						v.SetUint(num)
					default:
						log.Warn().Msgf("Not a support type '%s'", v.Kind())
					}
				}
			}
		}
	}
}

func fieldIsExported(field reflect.StructField) bool {
	return field.Name[0] >= 65 == true && field.Name[0] <= 90 == true
}

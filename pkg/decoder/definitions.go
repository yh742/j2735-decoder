package decoder

import (
	"strings"
)

// StringFormatType is used to identify which format to decode
type StringFormatType uint8

// String formats that the module supports
const (
	XML StringFormatType = iota
	JSON
	PASS
)

// UnmarshalYAML for decoder.StringFormatType type
func (format *StringFormatType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var formatStr string
	if err := unmarshal(&formatStr); err != nil {
		return err
	}
	switch strings.ToLower(formatStr) {
	case "xml":
		*format = XML
	case "json":
		*format = JSON
	case "pass":
		*format = PASS
	}
	return nil
}

// ParseString returns number based on string representation
func (format *StringFormatType) ParseString(str string) (uint8, bool) {
	m := map[string]uint8{"xml": 0, "json": 1, "pass": 2}
	val, ok := m[str]
	return val, ok
}

// ID to identify message type
const (
	SPaT int64 = 19
	BSM  int64 = 20
	EVA  int64 = 22
	RSA  int64 = 27
	PSM  int64 = 32
)

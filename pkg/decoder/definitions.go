package decoder

import (
	"fmt"
	"strings"
)

// StringFormatType is used to identify which format to decode
type StringFormatType uint8

// String formats that the module supports
const (
	XML StringFormatType = iota
	JSON
	PASS
	NA
)

// MarshalJSON for decoder.StringFormatType type
func (format StringFormatType) MarshalJSON() ([]byte, error) {
	switch format {
	case XML:
		return []byte(`"xml"`), nil
	case JSON:
		return []byte(`"json"`), nil
	case PASS:
		return []byte(`"pass"`), nil
	}
	return nil, fmt.Errorf("could not marshal format type to json %v", format)
}

// UnmarshalJSON for decoder.StringFormatType type
func (format *StringFormatType) UnmarshalJSON(b []byte) error {
	str := string(b)
	temp, err := parseFormatType(strings.ReplaceAll(str, "\"", ""))
	if err != nil {
		return err
	}
	*format = temp
	return nil
}

// UnmarshalYAML for decoder.StringFormatType type
func (format *StringFormatType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var formatStr string
	if err := unmarshal(&formatStr); err != nil {
		return err
	}
	temp, err := parseFormatType(formatStr)
	if err != nil {
		return err
	}
	*format = temp
	return nil
}

// ParseFormatType is a helper method for parsing format types
func parseFormatType(formatStr string) (StringFormatType, error) {
	switch strings.ToLower(formatStr) {
	case "xml":
		return XML, nil
	case "json":
		return JSON, nil
	case "pass":
		return PASS, nil
	}
	return NA, fmt.Errorf("cannot convert %s", formatStr)
}

// ParseString returns number based on string representation
func (format *StringFormatType) ParseString(str string) (uint8, bool) {
	formatType, err := parseFormatType(str)
	if err != nil {
		return 0, false
	}
	return uint8(formatType), true
}

// ID to identify message type
const (
	SPaT int64 = 19
	BSM  int64 = 20
	EVA  int64 = 22
	RSA  int64 = 27
	PSM  int64 = 32
)

package decoder

// StringFormatType is used to identify which format to decode
type StringFormatType int

// String formats that the module supports
const (
	XML  StringFormatType = iota
	JSON StringFormatType = iota
)

// ID to identify message type
const (
	SPaT int64 = 19
	BSM  int64 = 20
	EVA  int64 = 22
	RSA  int64 = 27
	PSM  int64 = 32
)
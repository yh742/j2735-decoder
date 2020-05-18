package decoder

// #cgo CFLAGS: -I${SRCDIR}/c/
// #cgo LDFLAGS: -L${SRCDIR}/c/ -lasncodec
// #include <MessageFrame.h>
// #include <xer_encoder.h>
// #include <per_decoder.h>
// void free_struct(asn_TYPE_descriptor_t descriptor, void* frame) {
// 		ASN_STRUCT_FREE(descriptor, frame);
// }
import "C"
import (
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
)

// DecodeBytes is a public function for other packages to decode string.
// It returns a string in either json or xml format.
// If source is provided, the broker topic will be added.
// Use pb dictate whether to deserialize first
func DecodeBytes(bytes []byte, format StringFormatType, source string, pb bool) (string, error) {
	var msgFrame *C.MessageFrame_t = nil
	// deserialize if (1) spat message (2) protobuf
	if pb && strings.Contains(strings.ToLower(source), "spat") {
		var clientMsg ClientDestinedMsg
		if err := proto.Unmarshal(bytes, &clientMsg); err != nil {
			log.Error().Msg("Cannot deserialize protobuf message")
			return "", errors.New("Cannot deserialize protobuf message")
		}
		cliBytes := clientMsg.GetMsgBytes()
		msgFrame = decodeMessageFrame(&C.asn_DEF_MessageFrame, cliBytes, uint64(len(cliBytes)))
	} else {
		msgFrame = decodeMessageFrame(&C.asn_DEF_MessageFrame, bytes, uint64(len(bytes)))
	}
	if msgFrame == nil {
		log.Error().
			Msg("Cannot decode bytes to messageframe struct")
		return "", errors.New("Cannot decode bytes to messageframe struct")
	}
	defer C.free_struct(C.asn_DEF_MessageFrame, unsafe.Pointer(msgFrame))
	log.Debug().
		Msgf("Decoding message type: %d", int64(msgFrame.messageId))
	// decode in different formats
	switch format {
	case JSON:
		xml, err := msgFrameToXMLString(msgFrame)
		if err != nil {
			return "", errors.New("decoding json error")
		}
		if source != "" {
			xml = strings.Replace(xml, "<MessageFrame>", "<MessageFrame><source>"+source+"</source>", 1)
		}
		return xmlStringToJSONString(xml)
	case XML:
		xml, err := msgFrameToXMLString(msgFrame)
		if err != nil {
			return "", errors.New("decoding xml error")
		}
		if source != "" {
			xml = strings.Replace(xml, "<MessageFrame>", "<MessageFrame><source>"+source+"</source>", 1)
		}
		return xml, nil
	default:
		return "", errors.New("format type not supported")
	}
}

// decodeMessageFrame requires caller to free the MessageFrame returned
func decodeMessageFrame(descriptor *C.asn_TYPE_descriptor_t, bytes []byte, length uint64) *C.MessageFrame_t {
	var decoded unsafe.Pointer
	cBytes := C.CBytes(bytes)
	defer C.free(cBytes)
	rval := C.uper_decode_complete(
		nil,
		descriptor,
		&decoded,
		cBytes,
		C.ulong(length))
	if rval.code != C.RC_OK {
		err := fmt.Sprintf("Broken Rectangle encoding at byte %d", (uint64)(rval.consumed))
		log.Error().
			Msg(err)
		return nil
	}
	return (*C.MessageFrame_t)(decoded)
}

package decoder

// #cgo CFLAGS: -I${SRCDIR}/c/
// #cgo LDFLAGS: -L${SRCDIR}/c/ -lasncodec
// #include <MessageFrame.h>
// #include <xer_encoder.h>
// #include <per_decoder.h>
// int xer__print2s (const void *buffer, size_t size, void *app_key)
// {
//     char *string = (char *) app_key;
//     strncat(string, buffer, size);
//     return 0;
// }
// int xer_sprint(void *string, asn_TYPE_descriptor_t *td, void *sptr)
// {
//     asn_enc_rval_t er;
//     er = xer_encode(td, sptr, XER_F_CANONICAL, xer__print2s, string);
//     if (er.encoded == -1)
//         return -1;
//     return er.encoded;
// }
// size_t partIIcontent_size()
// {
//		return sizeof(PartIIcontent_143P0_t);
// }
// struct PartIIcontent* get_partII(BasicSafetyMessage_t* ptr, int index)
// {
// 		return ptr->partII->list.array[index];
// }
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	xj "github.com/basgys/goxml2json"
	"github.com/rs/zerolog/log"
)

// octetStringToGoString takes in a ASN1 octet string and converts it to a Go string in hex
func octetStringToGoString(oString *C.OCTET_STRING_t) string {
	size := int(oString.size)
	str := ""
	for x := 0; x < size; x++ {
		octetByte := *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(oString.buf)) + uintptr(x)))
		str += fmt.Sprintf("%02X ", octetByte)
	}
	return str
}

// bitStringToGoString takes in a ASN1 bit string and converts it to a Go string in binary
func bitStringToGoString(bString *C.BIT_STRING_t) string {
	bitsUnused := uint64(bString.bits_unused)
	size := uint64(bString.size)
	body := uint8(*bString.buf)
	bits := int((size * 8) - bitsUnused)
	resStr := ""
	for x := 0; x < bits; x++ {
		if x := 0x80 & body; x == 128 {
			resStr += "1"
		} else {
			resStr += "0"
		}
		body = body << 1
	}
	return resStr
}

// msgFrameToXMLString convert message frame to XML
func msgFrameToXMLString(msgFrame *C.MessageFrame_t) (string, error) {
	size := 4096
	var buffer []byte
	for true {
		buffer = make([]byte, size)
		bufPtr := unsafe.Pointer(&buffer[0])
		rval := C.xer_sprint(bufPtr, &C.asn_DEF_MessageFrame, unsafe.Pointer(msgFrame))
		log.Debug().Msgf("Bytes Encoded: %d", int(rval))
		size = int(rval)
		if int(rval) == -1 {
			err := "Cannot encode message!"
			log.Error().Msg(err)
			return "", errors.New(err)
		} else if int(rval) > len(buffer) {
			continue
		}
		break
	}
	return fmt.Sprintf("%s", buffer[:size]), nil
}

// xmlStringToJsonString converts xml encoded string to json
func xmlStringToJSONString(xmlStr string) (string, error) {
	xml := strings.NewReader(xmlStr)
	json, err := xj.Convert(xml)
	if err != nil {
		log.Error().Msgf("Cannot encode to JSON: %s", err)
		return "", err
	}
	return json.String(), nil
}

// AreEqualJSON compares to json to see if they're equal
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

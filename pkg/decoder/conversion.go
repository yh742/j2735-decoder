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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	xj "github.com/basgys/goxml2json"
)

func getPtrSize() uint64 {
	return strconv.IntSize / 8
}

func getSeqByIdx(seq unsafe.Pointer, idx uint64) unsafe.Pointer {
	return unsafe.Pointer(uintptr(seq) + uintptr(idx * getPtrSize()))
}

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
	size := 2048
	var buffer []byte
	for true {
		buffer = make([]byte, size)
		bufPtr := unsafe.Pointer(&buffer[0])
		rval := C.xer_sprint(bufPtr, &C.asn_DEF_MessageFrame, unsafe.Pointer(msgFrame))
		Logger.Infof("Bytes Encoded: %d", int(rval))
		if int(rval) == -1 {
			err := "Cannot encode message!"
			Logger.Error(err)
			return "", errors.New(err)
		} else if int(rval) > len(buffer) {
			size = int(rval)
			continue
		}
		break
	}
	return fmt.Sprintf("%s", buffer), nil
}

// xmlStringToJsonString converts xml encoded string to json
func xmlStringToJSONString(xmlStr string) (string, error) {
	xml := strings.NewReader(xmlStr)
	json, err := xj.Convert(xml)
	if err != nil {
		Logger.Errorf("Cannot encode to JSON: %s", err)
		return "", err
	}
	return json.String(), nil
}

// msgFrameToSDMapBSM converts message frames to format ingested by SDMAP
func msgFrameToSDMapBSM(msgFrame *C.MessageFrame_t) (*MapAgtBSM, error) {
	if int64(msgFrame.messageId) != BSM {
		return nil, errors.New("this is not the right message type")
	}
	bsm := (*C.BasicSafetyMessage_t)(unsafe.Pointer(&msgFrame.value.choice))
	coreData := bsm.coreData
	partII := bsm.partII
	sdData := &MapAgtBSM{
		MsgCnt:  int64(coreData.msgCnt),
		Lat:     int64(coreData.lat),
		Long:    int64(coreData.Long),
		Elev:    int64(coreData.elev),
		Speed:   int64(coreData.speed),
		Heading: int64(coreData.heading),
		Angle:   int64(coreData.angle),
		EV:      int64(0),
	}
	sdData.ID = strings.TrimSpace(octetStringToGoString(&coreData.id))
	if partII != nil {
		const PtrSize = strconv.IntSize / 8
		for i := uint64(0); i < uint64(partII.list.count); i++ {
			contentPtr := (**C.PartIIcontent_143P0_t)(getSeqByIdx(unsafe.Pointer(partII.list.array), i))
			switch uint((*contentPtr).partII_Id) {
			// vehicle safety extension
			case 0:
				break
			// special vehicle extension
			case 1:
				specialVehicleExtensions := 
					(*C.SpecialVehicleExtensions_t)(unsafe.Pointer(&(*contentPtr).partII_Value.choice))
				if specialVehicleExtensions.vehicleAlerts != nil {
					sdData.EV = int64(specialVehicleExtensions.vehicleAlerts.sirenUse)
				}
				break
			// supplmental vehicle extension
			case 2:
				break
			// nothing is there or corrupt frames
			default:
				break
			}
		}
	}
	return sdData, nil
}

func numToPSMType(pType int64) string {
	switch pType {
	case 0:
		return "unavailable"
	case 1:
		return "aPEDESTRIAN"
	case 2:
		return "aPEDALCYCLIST"
	case 3:
		return "aPUBLICSAFETYWORKER"
	case 4:
		return "anANIMAL"
	default:
		return "unavailable"
	}
}

// msgFrameToSDMapPSM converts message frames to format ingested by SDMAP
func msgFrametoSDMapPSM(msgFrame *C.MessageFrame_t) (*MapAgtPSM, error) {
	if int64(msgFrame.messageId) != PSM {
		return nil, errors.New("this is not the right message type")
	}
	psmData := (*C.PersonalSafetyMessage_t)(unsafe.Pointer(&msgFrame.value.choice))
	sdData := &MapAgtPSM{
		MsgCnt:    int64(psmData.msgCnt),
		BasicType: numToPSMType(int64(psmData.basicType)),
		Lat:       int64(psmData.position.lat),
		Long:      int64(psmData.position.Long),
		Speed:     int64(psmData.speed),
		Heading:   int64(psmData.heading),
	}
	sdData.ID = strings.TrimSpace(octetStringToGoString(&psmData.id))
	return sdData, nil
}

// msgFrameToMapSPaT converts message frames to a SPaT format ingested by SDMAP
func msgFrametoMapSPaT(msgFrame *C.MessageFrame_t) (*SPaTList, error) {
	if int64(msgFrame.messageId) != SPaT {
		return nil, errors.New("this is not the right message type")
	}
	spatData := (*C.SPAT_t)(unsafe.Pointer(&msgFrame.value.choice))
	intersectionsCount := uint64(spatData.intersections.list.count)
	intersectionsPtr := unsafe.Pointer(spatData.intersections.list.array)
	var intersectionStates []IntersectionState
	for i := uint64(0); i < intersectionsCount; i++ {
		intersectionState := *(**C.IntersectionState_t)(getSeqByIdx(intersectionsPtr, i))
		id := fmt.Sprint(uint64(intersectionState.id.id))
		moy := uint64(*(intersectionState.moy))
		timeStamp := uint64(*(intersectionState.timeStamp))
		Logger.Debugf("IntersectionState: %d, id: %d, moy: %d, ts: %d", i, id, moy, timeStamp)
		movementStateCount := uint64((*intersectionState).states.list.count)
		movementStatePtr := unsafe.Pointer(intersectionState.states.list.array)
		var signalPhases []SignalPhaseGroup
		for i := uint64(0); i < movementStateCount; i++ { 
			signalPhaseGroup := *(**C.MovementState_t)(getSeqByIdx(movementStatePtr, i))
			signalGroupID := uint64(signalPhaseGroup.signalGroup)
			Logger.Debugf("SignalGroup: %d", signalGroupID)
			movementEventCount := uint64(signalPhaseGroup.state_time_speed.list.count)
			movementEventPtr := unsafe.Pointer(signalPhaseGroup.state_time_speed.list.array)
			for i := uint64(0); i < movementEventCount; i++ {
				movementEvent := *(**C.MovementEvent_t)(getSeqByIdx(movementEventPtr, i))
				state := uint64(movementEvent.eventState);
				minEnd := uint64(movementEvent.timing.minEndTime);
				maxEnd := uint64(*(movementEvent.timing.maxEndTime));
				Logger.Debugf("MovementEvent: %d, state: %d, minEnd: %d, maxEnd: %d", i, state, minEnd, maxEnd)
				sigPhase := SignalPhaseGroup {
					GroupID: signalGroupID,
					Status: state,
					MaxEndTime: maxEnd,
					MinEndTime: minEnd,
				}
				signalPhases = append(signalPhases, sigPhase)
			} 
		}
		state := IntersectionState {
			MinuteOfYear: moy,
			TimeStamp: timeStamp,
			SignalPhases: signalPhases,
		}
		state.ID = id
		intersectionStates = append(intersectionStates, state)
		Logger.Info("ADDED entry into interseciton states")
	}
	spatMsg := &SPaTList  {
		IntersectionStateList: intersectionStates,
	}
	return spatMsg, nil
}

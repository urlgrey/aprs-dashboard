package parser

// #cgo pkg-config: --libs libfap
/*
#include <fap.h>
#include <stdlib.h>

// type is a reserved keyword
fap_packet_type_t getPacketType(fap_packet_t* p) {
        return *p->type;
}
*/
import "C"
import (
	"errors"
	"strings"
	"time"
	"unsafe"

	"github.com/urlgrey/aprs-dashboard/models"
)

type AprsParser struct{}

func NewParser() *AprsParser {
	return &AprsParser{}
}

func (a *AprsParser) Initialize() {
	C.fap_init()
}

func (p *AprsParser) Close() {
	defer C.fap_cleanup()
}

func (p *AprsParser) ParseAprsPacket(rawPacket string, isAX25 bool) (message *models.AprsMessage, err error) {
	packet_cstring := C.CString(rawPacket)
	packet_length := C.uint(len(rawPacket))

	packet := C.fap_parseaprs(packet_cstring, packet_length, C.short(boolToInt(isAX25)))
	defer C.fap_free(packet)
	defer C.free(unsafe.Pointer(packet_cstring))

	if packet.error_code != nil {
		err = errors.New("Unable to parse APRS packet")
		return
	}

	message = &models.AprsMessage{
		Timestamp:           int32(time.Now().Unix()),
		SourceCallsign:      strings.ToUpper(C.GoString(packet.src_callsign)),
		DestinationCallsign: strings.ToUpper(C.GoString(packet.dst_callsign)),
		Latitude:            parseNilableFloat(packet.latitude),
		Longitude:           parseNilableFloat(packet.longitude),
		Speed:               parseNilableFloat(packet.speed),
		Course:              parseNilableUInt(packet.course),
		Altitude:            parseNilableFloat(packet.altitude),
		RawMessage:          C.GoStringN(packet.body, C.int(packet.body_len)),
	}

	switch C.getPacketType(packet) {
	case C.fapLOCATION:
		parsedMsg.PacketType = MessagePacketType
	case C.fapOBJECT:
		parsedMsg.PacketType = ObjectPacketType
	case C.fapITEM:
		parsedMsg.PacketType = ItemPacketType
	case C.fapMICE:
		parsedMsg.PacketType = MicePacketType
	case C.fapNMEA:
		parsedMsg.PacketType = NMEAPacketType
	case C.fapWX:
		parsedMsg.PacketType = WXPacketType
	case C.fapMESSAGE:
		parsedMsg.PacketType = MessagePacketType
	case C.fapCAPABILITIES:
		parsedMsg.PacketType = CapabilitiesPacketType
	case C.fapSTATUS:
		parsedMsg.PacketType = StatusPacketType
	case C.fapTELEMETRY:
		parsedMsg.PacketType = TelemetryPacketType
	case C.fapTELEMETRY_MESSAGE:
		parsedMsg.PacketType = TelemetryMessagePacketType
	case C.fapDX_SPOT:
		parsedMsg.PacketType = DXSpotPacketType
	case C.fapEXPERIMENTAL:
		parsedMsg.PacketType = ExperimentalPacketType
	}

	if packet.latitude != nil && packet.longitude != nil {
		message.IncludesPosition = true
	} else {
		message.IncludesPosition = false
	}

	if packet.wx_report != nil {
		message.Weather = &models.WeatherReport{
			Temperature:       parseNilableFloat(packet.wx_report.temp),
			InsideTemperature: parseNilableFloat(packet.wx_report.temp_in),
			Humidity:          parseNilableUInt(packet.wx_report.humidity),
			InsideHumidity:    parseNilableUInt(packet.wx_report.humidity_in),
			WindGust:          parseNilableFloat(packet.wx_report.wind_gust),
			WindDirection:     parseNilableUInt(packet.wx_report.wind_dir),
			WindSpeed:         parseNilableFloat(packet.wx_report.wind_speed),
		}
	}

	if message.SourceCallsign == "" {
		err = errors.New("Unable to find source callsign in APRS packet")
	}

	return
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func parseNilableFloat(d *C.double) float64 {
	if d != nil {
		return float64(C.double(*d))
	} else {
		return 0
	}
}

func parseNilableUInt(d *C.uint) uint8 {
	if d != nil {
		return uint8(C.uint(*d))
	} else {
		return 0
	}
}

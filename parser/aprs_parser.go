package parser

// #cgo pkg-config: --libs libfap
// #include <fap.h>
// #include <stdlib.h>
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
	C.fap_init()
	return &AprsParser{}
}

func (p *AprsParser) Finish() {
	defer C.fap_cleanup()
}

func (p *AprsParser) ParseAprsPacket(message string, isAX25 bool) (*models.AprsMessage, error) {
	message_cstring := C.CString(message)
	message_length := C.uint(len(message))

	packet := C.fap_parseaprs(message_cstring, message_length, C.short(boolToInt(isAX25)))
	defer C.fap_free(packet)
	defer C.free(unsafe.Pointer(message_cstring))

	if packet.error_code != nil {
		return &models.AprsMessage{}, errors.New("Unable to parse APRS message")
	}

	parsedMsg := models.AprsMessage{
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
	if packet.latitude != nil && packet.longitude != nil {
		parsedMsg.IncludesPosition = true
	} else {
		parsedMsg.IncludesPosition = false
	}

	if packet.wx_report != nil {
		w := models.WeatherReport{
			Temperature:       parseNilableFloat(packet.wx_report.temp),
			InsideTemperature: parseNilableFloat(packet.wx_report.temp_in),
			Humidity:          parseNilableUInt(packet.wx_report.humidity),
			InsideHumidity:    parseNilableUInt(packet.wx_report.humidity_in),
			WindGust:          parseNilableFloat(packet.wx_report.wind_gust),
			WindDirection:     parseNilableUInt(packet.wx_report.wind_dir),
			WindSpeed:         parseNilableFloat(packet.wx_report.wind_speed),
		}
		parsedMsg.Weather = &w
	}

	if parsedMsg.SourceCallsign == "" {
		return nil, errors.New("Unable to find source callsign in APRS packet")
	}

	return &parsedMsg, nil
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

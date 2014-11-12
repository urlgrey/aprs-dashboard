package main

// #cgo pkg-config: --libs libfap
// #include <fap.h>
// #include <stdlib.h>
import "C"
import "unsafe"
import "errors"
import "time"

type AprsMessage struct {
	Timestamp           int32          `json:"timestamp"`
	SourceCallsign      string         `json:"src_callsign"`
	DestinationCallsign string         `json:"dst_callsign"`
	Status              string         `json:"status"`
	Symbol              string         `json:"symbol"`
	Latitude            float64        `json:"latitude"`
	Longitude           float64        `json:"longitude"`
	IncludesPosition    bool           `json:"includes_position"`
	Altitude            float64        `json:"altitude"`
	Speed               float64        `json:"speed"`
	Course              uint8          `json:"course"`
	Weather             *WeatherReport `json:"weather_report"`
	RawMessage          string         `json:"raw_message"`
}

type WeatherReport struct {
	Temperature       float64 `json:"temp"`
	InsideTemperature float64 `json:"temp_in"`
	Humidity          uint8   `json:"humidity"`
	InsideHumidity    uint8   `json:"humidity_in"`
	WindGust          float64 `json:"wind_gust"`
	WindDirection     uint8   `json:"wind_dir"`
	WindSpeed         float64 `json:"wind_speed"`
}

type AprsParser struct{}

func NewParser() *AprsParser {
	C.fap_init()
	return &AprsParser{}
}

func (p *AprsParser) Finish() {
	defer C.fap_cleanup()
}

func (p *AprsParser) parseAprsPacket(message string, isAX25 bool) (*AprsMessage, error) {
	message_cstring := C.CString(message)
	message_length := C.uint(len(message))

	packet := C.fap_parseaprs(message_cstring, message_length, C.short(boolToInt(isAX25)))
	defer C.fap_free(packet)
	defer C.free(unsafe.Pointer(message_cstring))

	if packet.error_code != nil {
		return &AprsMessage{}, errors.New("Unable to parse APRS message")
	}

	parsedMsg := AprsMessage{
		Timestamp:           int32(time.Now().Unix()),
		SourceCallsign:      C.GoString(packet.src_callsign),
		DestinationCallsign: C.GoString(packet.dst_callsign),
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
		w := WeatherReport{
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

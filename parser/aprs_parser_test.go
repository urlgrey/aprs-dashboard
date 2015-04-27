package parser

import (
	"testing"
	"time"
)

func Test_ParseAprsNonAX25PacketWithLocation(t *testing.T) {
	timeBeforeTest := int32(time.Now().Unix())
	p := NewParser()
	defer p.Finish()

	msg, err := p.ParseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)
	timeAfterTest := int32(time.Now().Unix())

	if err != nil {
		t.Error("Error was unexpectedly non-nil", err)
	}
	if msg.Timestamp < timeBeforeTest || msg.Timestamp > timeAfterTest {
		t.Error("Timestamp was set incorrectly")
	}
	if msg.SourceCallsign != "K7SSW" {
		t.Error("Source callsign not parsed correctly", msg.SourceCallsign)
	}
	if msg.DestinationCallsign != "APRS" {
		t.Error("Destination callsign not parsed correctly", msg.DestinationCallsign)
	}
	if msg.Latitude != 47.720333333333336 {
		t.Error("Latitude incorrect", msg.Latitude)
	}
	if msg.Longitude != -122.3735 {
		t.Error("Longitude incorrect", msg.Longitude)
	}
}

func Test_ParseAprsNonAX25PacketWithLocationAndSpeed(t *testing.T) {
	timeBeforeTest := int32(time.Now().Unix())
	p := NewParser()
	defer p.Finish()

	msg, err := p.ParseAprsPacket("KK6OLB-1>S8RUQU,WIDE1-1,WIDE2-1,qAR,KI6TEV-1:`2I\"l6h>/]\"43}Santa Rosa CA USA=", false)
	timeAfterTest := int32(time.Now().Unix())

	if err != nil {
		t.Error("Error was unexpectedly non-nil", err)
	}
	if msg.Timestamp < timeBeforeTest || msg.Timestamp > timeAfterTest {
		t.Error("Timestamp was set incorrectly")
	}
	if msg.SourceCallsign != "KK6OLB-1" {
		t.Error("Source callsign not parsed correctly", msg.SourceCallsign)
	}
	if msg.DestinationCallsign != "S8RUQU" {
		t.Error("Destination callsign not parsed correctly", msg.DestinationCallsign)
	}
	if msg.Latitude != 38.41916666666667 {
		t.Error("Latitude incorrect", msg.Latitude)
	}
	if msg.Longitude != -122.751 {
		t.Error("Longitude incorrect", msg.Longitude)
	}
	if msg.Speed != 3.704 {
		t.Error("Speed incorrect", msg.Speed)
	}
	if msg.Course != 20 {
		t.Error("Course incorrect", msg.Course)
	}
	if msg.Altitude != 28 {
		t.Error("Altitude incorrect", msg.Altitude)
	}
}

func Test_ParseAprsNonAX25PacketWithLocationAndWeather(t *testing.T) {
	timeBeforeTest := int32(time.Now().Unix())
	p := NewParser()
	defer p.Finish()

	msg, err := p.ParseAprsPacket("DW6161>APRS,TCPXX*,qAX,CWOP-4:@101509z5004.48N/00645.00E_049/000g000t046r000p019P013h97b10123WeatherCatV123B16H31", false)
	timeAfterTest := int32(time.Now().Unix())

	if err != nil {
		t.Error("Error was unexpectedly non-nil", err)
	}
	if msg.Timestamp < timeBeforeTest || msg.Timestamp > timeAfterTest {
		t.Error("Timestamp was set incorrectly")
	}
	if msg.SourceCallsign != "DW6161" {
		t.Error("Source callsign not parsed correctly", msg.SourceCallsign)
	}
	if msg.DestinationCallsign != "APRS" {
		t.Error("Destination callsign not parsed correctly", msg.DestinationCallsign)
	}
	if msg.Latitude != 50.074666666666666 {
		t.Error("Latitude incorrect", msg.Latitude)
	}
	if msg.Longitude != 6.75 {
		t.Error("Longitude incorrect", msg.Longitude)
	}
	if msg.Weather.Temperature != 7.777777777777778 {
		t.Error("Temperature incorrect", msg.Weather.Temperature)
	}
	if msg.Weather.InsideTemperature != 0 {
		t.Error("Temperature incorrect", msg.Weather.InsideTemperature)
	}
	if msg.Weather.Humidity != 97 {
		t.Error("Humidity incorrect", msg.Weather.Humidity)
	}
	if msg.Weather.InsideHumidity != 0 {
		t.Error("Inside Humidity incorrect", msg.Weather.InsideHumidity)
	}
	if msg.Weather.WindGust != 0 {
		t.Error("Wind Gust incorrect", msg.Weather.WindGust)
	}
	if msg.Weather.WindDirection != 49 {
		t.Error("Wind Direction incorrect", msg.Weather.WindDirection)
	}
	if msg.Weather.WindSpeed != 0 {
		t.Error("Wind Speed incorrect", msg.Weather.WindSpeed)
	}
}

func Test_ParseAprsUnsupportedFormatPacket(t *testing.T) {
	p := NewParser()
	defer p.Finish()

	msg, err := p.ParseAprsPacket("ZS6EY>APN982,ZS0TRG*,WIDE3-2,qAS,ZS6EY-1:g {UIV32N}", false)
	if err == nil {
		t.Error("Error was unexpectedly nil", err)
	}
	if msg.SourceCallsign != "" {
		t.Error("Source callsign not parsed correctly", msg.SourceCallsign)
	}
	if msg.DestinationCallsign != "" {
		t.Error("Destination callsign not parsed correctly", msg.DestinationCallsign)
	}
	if msg.Latitude != 0 {
		t.Error("Latitude incorrect", msg.Latitude)
	}
	if msg.Longitude != 0 {
		t.Error("Longitude incorrect", msg.Longitude)
	}
}

func Benchmark_ParseValidAprsPacket(b *testing.B) {
	p := NewParser()
	defer p.Finish()

	for i := 0; i < b.N; i++ {
		p.ParseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)
	}
}

func Benchmark_ParseInvalidAprsPacket(b *testing.B) {
	p := NewParser()
	defer p.Finish()

	for i := 0; i < b.N; i++ {
		p.ParseAprsPacket("ZS6EY>APN982,ZS0TRG*,WIDE3-2,qAS,ZS6EY-1:g {UIV32N}", false)
	}
}

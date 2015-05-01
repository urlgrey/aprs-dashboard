package models

type PacketType int

const (
	LocationPacketType PacketType = iota
	ObjectPacketType
	ItemPacketType
	MicePacketType
	NMEAPacketType
	WXPacketType
	MessagePacketType
	CapabilitiesPacketType
	StatusPacketType
	TelemetryPacketType
	TelemetryMessagePacketType
	DXSpotPacketType
	ExperimentalPacketType
)

type AprsMessage struct {
	PacketType
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

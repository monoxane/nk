package nk

import (
	"net"
	"sync"
)

// Aliases

type Model = string
type TBusAddress = uint8
type Level = uint32

// Structs

type Router struct {
	IP           string
	Address      TBusAddress
	Destinations uint16
	Sources      uint16
	Level        Level
	Matrix       Matrix
	Conn         net.Conn
	onUpdate     func(*Update)
}

type Destination struct {
	Label  string  `json:"label"`
	Id     uint16  `json:"id"`
	Source *Source `json:"source"`
}

type Source struct {
	Label string `json:"label"`
	Id    uint16 `json:"id"`
}

type nkRoutePacketPayload struct {
	NK2Header   uint32
	RTRAddress  TBusAddress
	UNKNB       uint16
	Destination uint16
	Source      uint16
	LevelMask   Level
	UNKNC       uint8
}

type nkRoutePacket struct {
	HeaderA uint32
	HeaderB uint16
	Payload nkRoutePacketPayload
	CRC     uint16
}

type CrosspointRequest struct {
	Source      uint16
	Destination uint16
	Level       Level
	Address     TBusAddress
}

type Matrix struct {
	destinations map[uint16]*Destination
	sources      map[uint16]*Source
	mux          sync.Mutex
}

type Update struct {
	Type string
	Data interface{}
}

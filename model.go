package nk

import "net"

// Aliases

type Model = string
type TBusAddress = uint8
type Destination = uint16
type Source = uint16
type Level = uint32

// Structs

type Router struct {
	IP           string
	Address      TBusAddress
	Destinations uint16
	Sources      uint16
	Level        Level
	Conn         net.Conn
}

type TBusPacketPayload struct {
	NK2Header   uint32
	RTRAddress  TBusAddress
	UNKNB       uint16
	Destination Destination
	Source      Source
	LevelMask   Level
	UNKNC       uint8
}

type TBusPacket struct {
	HeaderA uint32
	HeaderB uint16
	Payload TBusPacketPayload
	CRC     uint16
}

type CrosspointRequest struct {
	Source      Source
	Destination Destination
	Level       Level
	Address     TBusAddress
}

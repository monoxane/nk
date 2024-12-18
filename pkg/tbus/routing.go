package tbus

import (
	"bytes"
	"encoding/binary"

	"github.com/monoxane/nk"
	"github.com/monoxane/nk/pkg/tbus/crc"
)

type CrosspointRequest struct {
	Source      uint16
	Destination uint16
	Level       Level
	Address     TBusAddress
}

// GenerateXPTRequest Just returns payload to send to router to close xpt
func (xpt *CrosspointRequest) Packet() ([]byte, error) {
	destination := xpt.Destination - 1
	source := xpt.Source - 1

	payload := nkRoutePacketPayload{
		NK2Header:   0x4e4b3200,
		RTRAddress:  xpt.Address,
		UNKNB:       0x0409,
		Destination: destination,
		Source:      source,
		LevelMask:   xpt.Level,
		UNKNC:       0x00,
	}

	payloadBuffer := new(bytes.Buffer)
	err := binary.Write(payloadBuffer, binary.BigEndian, payload)
	if err != nil {
		nk.Log.Error().Err(err).Msg("unable pack tbus message into packet")
	}

	packet := nkRoutePacket{
		HeaderA: 0x50415332,
		HeaderB: 0x0012,
		Payload: payload,
		CRC:     crc.CRC16(payloadBuffer.Bytes()),
	}

	packetBuffer := new(bytes.Buffer)
	err = binary.Write(packetBuffer, binary.BigEndian, packet)
	if err != nil {
		nk.Log.Error().Err(err).Msg("unable to write packet to tbus gateway")
	}

	return packetBuffer.Bytes(), nil
}

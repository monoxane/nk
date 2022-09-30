package nk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func (rtr *Router) Route(dst uint16, src uint16) error {
	if dst <= 0 || dst > rtr.Destinations {
		return fmt.Errorf("requested destination is outside the range available on this router model")
	}

	if src <= 0 || src > rtr.Sources {
		return fmt.Errorf("requested source is outside the range available on this router model")
	}

	xptreq := CrosspointRequest{
		Source:      src,
		Destination: dst,
		Level:       rtr.Level,
		Address:     rtr.Address,
	}

	packet, err := xptreq.Packet()
	if err != nil {
		return errors.Wrap(err, "unable to generate crosspoint route request")
	}

	_, err = rtr.Conn.Write(packet)
	if err != nil {
		return errors.Wrap(err, "unable to send route request to router")
	}

	return nil
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
		log.Println("TBusPacketPayload binary.Write failed:", err)
	}

	packet := nkRoutePacket{
		HeaderA: 0x50415332,
		HeaderB: 0x0012,
		Payload: payload,
		CRC:     crc16(payloadBuffer.Bytes()),
	}

	packetBuffer := new(bytes.Buffer)
	err = binary.Write(packetBuffer, binary.BigEndian, packet)
	if err != nil {
		log.Println("TBustPacket binary.Write failed:", err)
	}

	return packetBuffer.Bytes(), nil
}

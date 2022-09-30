package nk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	LVL_MDVID uint32 = 1

	MODEL_NK_3G72      Model = "NK-3G72"
	MODEL_NK_3G64      Model = "NK-3G64"
	MODEL_NK_3G34      Model = "NK-3G34"
	MODEL_NK_3G16      Model = "NK-3G16"
	MODEL_NK_3G16_RCP  Model = "NK-3G16-RCP"
	MODEL_NK_3G164     Model = "NK-3G164"
	MODEL_NK_3G164_RCP Model = "NK-3G164-RCP"
)

var (
	NK2_KEEPALIVE        = []byte("HI")
	NK2_CONNECT_REQ      = []byte{0x50, 0x48, 0x4f, 0x45, 0x4e, 0x49, 0x58, 0x2d, 0x44, 0x42, 0x20, 0x4e, 0x0a}
	NK2_CONNECT_RESP     = []byte{0x57, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x0a}
	NK2_HEADER           = []byte{0x4e, 0x4b, 0x32}
	NK_STATUS_RESP       = []byte{0x05, 0x0B}
	NK_MULTI_STATUS_REQ  = []byte{0x50, 0x41, 0x53, 0x32, 0x00, 0x11, 0x4e, 0x4b, 0x32, 0x00, 0xfe, 0x02, 0x08, 0x00, 0x00, 0x00, 0x47, 0xff, 0xff, 0xff, 0xff, 0xc7, 0x08}
	NK_MULTI_STATUS_RESP = []byte{0x03, 0xe1}
)

func New(IP string, RTRAddress uint8, model Model) *Router {
	rtr := &Router{
		IP:      IP,
		Address: RTRAddress,
	}

	switch model {
	case MODEL_NK_3G72:
		rtr.Destinations = 72
		rtr.Sources = 72
		rtr.Level = LVL_MDVID
	case MODEL_NK_3G64:
		rtr.Destinations = 64
		rtr.Sources = 64
		rtr.Level = LVL_MDVID
	case MODEL_NK_3G34:
		rtr.Destinations = 34
		rtr.Sources = 34
		rtr.Level = LVL_MDVID
	case MODEL_NK_3G16, MODEL_NK_3G16_RCP:
		rtr.Destinations = 16
		rtr.Sources = 16
		rtr.Level = LVL_MDVID
	case MODEL_NK_3G164, MODEL_NK_3G164_RCP:
		rtr.Destinations = 4
		rtr.Sources = 16
		rtr.Level = LVL_MDVID
	}

	rtr.Matrix.destinations = make(map[uint16]*Destination)
	rtr.Matrix.sources = make(map[uint16]*Source)

	for i := 0; i < int(rtr.Sources)+1; i++ {
		rtr.Matrix.sources[uint16(i)] = &Source{
			Id:    uint16(i),
			Label: fmt.Sprintf("IN %d", i),
		}
	}

	for i := 0; i < int(rtr.Destinations)+1; i++ {
		rtr.Matrix.destinations[uint16(i)] = &Destination{
			Id:    uint16(i),
			Label: fmt.Sprintf("OUT %d", i),
		}
	}

	return rtr
}

func (rtr *Router) LoadLabels(labels string) {
	lines := strings.Split(labels, "\n")
	for i, line := range lines {
		columns := strings.Split(line, ",")
		if len(columns) < 4 {
			continue
		}

		log.Printf("%+v", columns)
		if _, ok := rtr.Matrix.destinations[uint16(i+1)]; ok {
			rtr.Matrix.destinations[uint16(i+1)].Setlabel(columns[1])
		}
		if _, ok := rtr.Matrix.sources[uint16(i+1)]; ok {
			rtr.Matrix.sources[uint16(i+1)].Setlabel(columns[3])
		}
	}
}

func (rtr *Router) Connect() error {
	conn, err := net.Dial("tcp", rtr.IP+":5000")
	if err != nil {
		log.Fatalln(err)
	}
	rtr.Conn = conn
	defer rtr.Conn.Close()

	if _, err = rtr.Conn.Write(NK2_CONNECT_REQ); err != nil {
		log.Printf("failed to send the client request: %v\n", err)
	}

	go func() {
		for range time.Tick(10 * time.Second) {
			rtr.Conn.Write([]byte("HI"))
		}
	}()

	for {
		buf := make([]byte, 2048)
		len, err := rtr.Conn.Read(buf)
		switch err {
		case nil:
			rtr.processNKMessage(buf, len)
		case io.EOF:
			return errors.New("remote connection closed")
		default:
			return errors.Wrap(err, "unhandled server error")
		}
	}
}

func (rtr *Router) processNKMessage(buffer []byte, length int) {
	msg := buffer[:length]
	log.Printf("Processing message of len %d: %x", length, msg)

	if length == len(NK2_CONNECT_RESP) && bytes.Equal(msg, NK2_CONNECT_RESP) {
		log.Printf("Sucessfully Connected")
		rtr.Conn.Write(NK_MULTI_STATUS_REQ)
	}

	if length > 3 && bytes.Equal(msg[:3], NK2_HEADER) {
		log.Printf("NK Command or Response, CMD: %x", msg[5:7])
		if bytes.Equal(msg[5:7], NK_STATUS_RESP) {
			rtr.parseSingleUpdateMessage(msg)
		}

		if bytes.Equal(msg[5:7], NK_MULTI_STATUS_RESP) {
			rtr.parseMultiUpdateMessage(msg)
		}
	}
}

func (rtr *Router) parseSingleUpdateMessage(msg []byte) {
	dst := binary.BigEndian.Uint16(msg[8:10]) + 1
	src := binary.BigEndian.Uint16(msg[10:12]) + 1
	lvl := binary.BigEndian.Uint32(msg[12:16])

	rtr.updateMatrix(lvl, dst, src)
}

func (rtr *Router) parseMultiUpdateMessage(msg []byte) {
	table := msg[15 : len(msg)-2]

	currentCrosspointByte := 1
	for {
		if currentCrosspointByte >= len(table) {
			break
		}

		dst := uint16(currentCrosspointByte/3) + 1
		src := binary.BigEndian.Uint16(table[currentCrosspointByte:currentCrosspointByte+2]) + 1
		lvl := uint32(1)

		rtr.updateMatrix(lvl, dst, src)

		currentCrosspointByte++
		currentCrosspointByte++
		currentCrosspointByte++
	}

}

func (rtr *Router) updateMatrix(lvl Level, dst uint16, src uint16) {
	if lvl == rtr.Level {
		log.Printf("Updating Crosspoint State: DST %2d SRC %2d", dst, src)
		rtr.Matrix.SetCrosspoint(dst, src)
	}
}

// 4e 4b 32 00 02 03 e1 0000004701000
// 4e 4b 32 00 02 03 e1 00 00 00 47 01 00 00 00 01
// 0034ff
// 0035ff
// 0037ff
// 0036ff
// 003aff
// 003bff
// 0020ff
// 001fff
// 0021ff
// 0022ff
// 002aff
// 002bff
// 0023ff
// 0024ff
// 0031ff
// 0041ff
// 0042ff
// 0043ff
// 0025ff
// 0042ff
// 002eff
// 0032ff
// 0033ff
// 0021ff
// 003cff
// 003cff
// 0004ff
// 0005ff
// 002fff
// 0030ff
// 002dff
// 0010ff
// 0007ff
// 001dff
// 0016ff
// 0016ff
// 002eff
// 0033ff
// 003dff
// 003dff
// 0002ff
// 0003ff
// 000eff
// 000fff
// 001dff
// 001eff
// 0011ff
// 0006ff
// 003aff
// 003bff
// 000fff
// 003dff
// 0002ff
// 0003ff
// 000fff
// 003dff
// 0017ff
// 0001ff
// 0031ff
// 0016ff
// 0000ff
// 0018ff
// 000cff
// 0016ff
// 0036ff
// 0045ff
// 0047ff
// 0047ff
// 0000ff
// 000eff
// 0001ff
// 0010ff
// cdbe

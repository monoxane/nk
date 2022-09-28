package nk

import (
	"bufio"
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

	return rtr
}

func (rtr *Router) Connect() {
	conn, err := net.Dial("tcp", rtr.IP+":5000")
	if err != nil {
		log.Fatalln(err)
	}
	rtr.Conn = conn
	defer rtr.Conn.Close()

	serverReader := bufio.NewReader(rtr.Conn)

	openConnStr := strings.TrimSpace("PHEONIX-DB")
	if _, err = rtr.Conn.Write([]byte(openConnStr + "\n")); err != nil {
		log.Printf("failed to send the client request: %v\n", err)
	}

	go func() {
		for range time.Tick(10 * time.Second) {
			rtr.Conn.Write([]byte("HI"))
		}
	}()

	for {

		serverResponse, err := serverReader.ReadString('\n')
		switch err {
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}
}

func (rtr *Router) SetCrosspoint(level Level, destination Destination, source Source) error {
	if level != rtr.Level {
		return fmt.Errorf("requested level is not possible on this router")
	}

	if destination <= 0 || destination > rtr.Destinations {
		return fmt.Errorf("requested destination is outside the range available on this router model")
	}

	if source <= 0 || source > rtr.Sources {
		return fmt.Errorf("requested source is outside the range available on this router model")
	}

	xptreq := CrosspointRequest{
		Source:      source,
		Destination: destination,
		Level:       level,
		Address:     rtr.Address,
	}

	packet, err := xptreq.GeneratePacket()
	if err != nil {
		return errors.Wrap(err, "unable to generate crosspoint route request")
	}

	_, err = rtr.Conn.Write(packet)
	if err != nil {
		return errors.Wrap(err, "unable to send route request to router")
	}

	return nil
}

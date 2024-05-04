package router

import (
	"fmt"
	"net"

	"github.com/monoxane/nk/pkg/levels"
	"github.com/monoxane/nk/pkg/matrix"
	"github.com/monoxane/nk/pkg/models"
	"github.com/monoxane/nk/pkg/tbus"
)

type RouteUpdate struct {
	Type        string
	Destination *matrix.Destination
	Source      *matrix.Source
}

// A Router is a fully featured NK Routing matrix with state and label management
type Router struct {
	// The NK-IPS or NK-NET interface used to access this Router
	gateway *tbus.TBusGateway

	// The metadata of this router
	Address      tbus.TBusAddress
	Destinations uint16
	Sources      uint16
	Level        tbus.Level
	Matrix       matrix.Matrix

	// client facing update messages
	onRouteUpdate  func(*RouteUpdate)
	onStatusUpdate func(tbus.StatusUpdate)
}

func New(ip net.IP, routerAddress tbus.TBusAddress, model models.Model) *Router {
	rtr := &Router{}

	switch model {
	case models.NK_3G72:
		rtr.Destinations = 72
		rtr.Sources = 72
		rtr.Level = levels.MD_Vid
	case models.NK_3G64:
		rtr.Destinations = 64
		rtr.Sources = 64
		rtr.Level = levels.MD_Vid
	case models.NK_3G34:
		rtr.Destinations = 34
		rtr.Sources = 34
		rtr.Level = levels.MD_Vid
	case models.NK_3G16, models.NK_3G16_RCP:
		rtr.Destinations = 16
		rtr.Sources = 16
		rtr.Level = levels.MD_Vid
	case models.NK_3G164, models.NK_3G164_RCP:
		rtr.Destinations = 4
		rtr.Sources = 16
		rtr.Level = levels.MD_Vid
	}

	rtr.Matrix.Init(rtr.Destinations, rtr.Sources)

	gw := tbus.NewGateway(ip, rtr.handleRouteUpdate, rtr.handleStatusUpdate)

	rtr.gateway = gw

	return rtr
}

func (rtr *Router) Connect() error {
	return rtr.gateway.Connect()
}

func (rtr *Router) Disconnect() {
	rtr.gateway.Disconnect()
}

func (rtr *Router) SetOnUpdate(notify func(*RouteUpdate)) {
	rtr.onRouteUpdate = notify
}

func (rtr *Router) Route(dst uint16, src uint16) error {
	if dst <= 0 || dst > rtr.Destinations {
		return fmt.Errorf("requested destination is outside the range available on this router model")
	}

	if src <= 0 || src > rtr.Sources {
		return fmt.Errorf("requested source is outside the range available on this router model")
	}

	return rtr.gateway.Route(rtr.Address, rtr.Level, dst, src)
}

func (rtr *Router) handleRouteUpdate(update tbus.RouteUpdate) {
	rtr.updateMatrix(uint16(update.Destination), uint16(update.Source))
}

func (rtr *Router) handleStatusUpdate(update tbus.StatusUpdate) {
	if rtr.onStatusUpdate != nil {
		rtr.onStatusUpdate(update)
	}
}

func (rtr *Router) updateMatrix(dst uint16, src uint16) {
	rtr.Matrix.SetCrosspoint(dst, src)

	if rtr.onRouteUpdate != nil {
		go rtr.onRouteUpdate(&RouteUpdate{
			Type:        "destination",
			Destination: rtr.Matrix.GetDestination(dst),
		})
	}
}

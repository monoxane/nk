package router

import (
	"strings"

	"github.com/monoxane/nk/pkg/matrix"
)

func (rtr *Router) LoadLabels(labels string) {
	lines := strings.Split(labels, "\n")
	for i, line := range lines {
		columns := strings.Split(line, ",")
		if len(columns) < 4 {
			continue
		}

		if i < int(rtr.Destinations) {
			rtr.Matrix.GetDestination(uint16(i + 1)).SetLabel(columns[1])
		}

		if i < int(rtr.Sources) {
			rtr.Matrix.GetSource(uint16(i + 1)).SetLabel(columns[3])
		}
	}
}

func (rtr *Router) UpdateSourceLabel(src int, label string) {
	if src <= int(rtr.Sources) {
		rtr.Matrix.GetSource(uint16(src)).SetLabel(label)
		go rtr.onRouteUpdate(&RouteUpdate{
			Type:   "source",
			Source: rtr.Matrix.GetSource(uint16(src)),
		})

		rtr.Matrix.ForEachDestination(func(i uint16, dst *matrix.Destination) {
			if dst.Source != nil && dst.Source.GetID() == uint16(src) {
				if rtr.onRouteUpdate != nil {
					go rtr.onRouteUpdate(&RouteUpdate{
						Type:        "destination",
						Destination: dst,
					})
				}
			}
		})
	}
}

func (rtr *Router) UpdateDestinationLabel(dst int, label string) {
	if dst <= int(rtr.Destinations) {
		rtr.Matrix.GetDestination(uint16(dst)).SetLabel(label)

		if rtr.onRouteUpdate != nil {
			go rtr.onRouteUpdate(&RouteUpdate{
				Type:        "destination",
				Destination: rtr.Matrix.GetDestination(uint16(dst)),
			})
		}
	}
}

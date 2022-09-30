package nk

import "log"

func (matrix *Matrix) SetCrosspoint(dst uint16, src uint16) {
	matrix.mux.Lock()
	defer matrix.mux.Unlock()
	matrix.destinations[dst].Source = matrix.sources[src]
	log.Printf("%+v %+v", matrix.destinations[dst], matrix.sources[src])
}

func (matrix *Matrix) GetDestination(dst uint16) *Destination {
	matrix.mux.Lock()
	defer matrix.mux.Unlock()
	return matrix.destinations[dst]
}

func (dst *Destination) GetID() uint16 {
	return dst.Id
}

func (dst *Destination) GetLabel() string {
	return dst.Label
}

func (dst *Destination) Setlabel(lbl string) {
	dst.Label = lbl
}

func (dst *Destination) GetSource() *Source {
	return dst.Source
}

func (src *Source) GetID() uint16 {
	return src.Id
}

func (src *Source) GetLabel() string {
	return src.Label
}

func (src *Source) Setlabel(lbl string) {
	src.Label = lbl
}

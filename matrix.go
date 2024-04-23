package nk

import (
	"encoding/json"
	"sort"
)

func (matrix *Matrix) MarshalJSON() ([]byte, error) {
	type res struct {
		Destinations []*Destination `json:"destinations,omitempty"`
		Sources      []*Source      `json:"sources,omitempty"`
	}

	var result res

	for _, dst := range matrix.destinations {
		if dst != nil && dst.Id != 0 {
			result.Destinations = append(result.Destinations, dst)
		}
	}

	sort.Slice(result.Destinations, func(i, j int) bool {
		return result.Destinations[i].Id < result.Destinations[j].Id
	})

	for _, src := range matrix.sources {
		if src != nil && src.Id != 0 {
			result.Sources = append(result.Sources, src)
		}
	}

	sort.Slice(result.Sources, func(i, j int) bool {
		return result.Sources[i].Id < result.Sources[j].Id
	})

	return json.Marshal(&result)
}

func (matrix *Matrix) SetCrosspoint(dst uint16, src uint16) {
	matrix.mux.Lock()
	defer matrix.mux.Unlock()
	matrix.destinations[dst].Source = matrix.sources[src]
}

func (matrix *Matrix) GetDestination(dst uint16) *Destination {
	matrix.mux.Lock()
	defer matrix.mux.Unlock()
	return matrix.destinations[dst]
}

func (matrix *Matrix) GetSource(src uint16) *Source {
	matrix.mux.Lock()
	defer matrix.mux.Unlock()
	return matrix.sources[src]
}

func (dst *Destination) GetID() uint16 {
	return dst.Id
}

func (dst *Destination) GetIDInt() int {
	return int(dst.Id)
}

func (dst *Destination) GetLabel() string {
	return dst.Label
}

func (dst *Destination) SetLabel(lbl string) {
	dst.Label = lbl
}

func (dst *Destination) GetSource() *Source {
	return dst.Source
}

func (src *Source) GetID() uint16 {
	return src.Id
}

func (src *Source) GetIDInt() int {
	return int(src.Id)
}

func (src *Source) GetLabel() string {
	return src.Label
}

func (src *Source) SetLabel(lbl string) {
	src.Label = lbl
}

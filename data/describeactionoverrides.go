package data

import (
	"github.com/goforce/api/soap"
	"strings"
)

type ActionOverrides struct {
	byName map[string]*ActionOverride
}

func newActionOverrides(aos []*soap.ActionOverride) *ActionOverrides {
	m := make(map[string]*ActionOverride)
	for _, ao := range aos {
		m[strings.ToLower(ao.Name)] = &ActionOverride{ao}
	}
	return &ActionOverrides{m}
}

func (aos *ActionOverrides) List() []*ActionOverride {
	l := make([]*ActionOverride, 0, len(aos.byName))
	for _, ao := range aos.byName {
		l = append(l, ao)
	}
	return l
}

func (aos *ActionOverrides) Get(name string) (*ActionOverride, bool) {
	v, ok := aos.byName[strings.ToLower(name)]
	return v, ok
}

type ActionOverride struct {
	dr *soap.ActionOverride
}

func (ao *ActionOverride) IsAvailableInTouch() bool { return ao.dr.IsAvailableInTouch }
func (ao *ActionOverride) Name() string             { return ao.dr.Name }
func (ao *ActionOverride) PageId() string           { return ao.dr.PageId }
func (ao *ActionOverride) Url() string              { return ao.dr.Url }

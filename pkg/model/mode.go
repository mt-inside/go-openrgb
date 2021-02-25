package model

import (
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/wire"
)

type Mode interface {
	GetName() string
	render(indent int) []indentedString
}

type EffectMode struct {
	name      string
	flags     wire.ModeFlags
	minSpeed  uint32
	speed     uint32
	maxSpeed  uint32
	direction wire.EffectDirection
	colorMode wire.ColorMode
	minColors uint32
	maxColors uint32
	colors    []*colorful.Color
}

func (m *EffectMode) GetName() string {
	return m.name
}

func (m *EffectMode) render(indent int) []indentedString {
	attrs := []string{}
	if !wire.ModeDirectionFlagsClear(m.flags) {
		attrs = append(attrs, fmt.Sprintf("direction: %s", m.direction))
	}
	if m.flags&wire.HasSpeed == wire.HasSpeed {
		attrs = append(attrs, fmt.Sprintf("speed: %d (%d-%d)", m.speed, m.minSpeed, m.maxSpeed))
	}

	ss := []indentedString{
		{indent, fmt.Sprintf("MODE [Effect] %s (%s)", m.name, m.flags)},
		{indent + 1, strings.Join(attrs, ",")},
	}
	if m.colorMode == wire.ColorModeNone {
		ss = append(ss, indentedString{indent + 1, fmt.Sprintf("Colors: %s", m.colorMode)})
	} else {
		ss = append(ss, indentedString{indent + 1, fmt.Sprintf("Colors: %s, (%d-%d) %s", m.colorMode, m.minColors, m.maxColors, renderColors(m.colors))})
	}

	return ss
}
func (m *EffectMode) String() string {
	return renderIndents(m.render(0))
}

type DirectMode struct {
	name  string
	Zones []*Zone
}

func (m *DirectMode) GetName() string {
	return m.name
}

func (m *DirectMode) render(indent int) []indentedString {
	ss := []indentedString{
		{indent, fmt.Sprintf("MODE [Direct] %s", m.name)},
	}

	for _, zone := range m.Zones {
		ss = append(ss, zone.render(indent+1)...)
	}

	return ss
}
func (m *DirectMode) String() string {
	return renderIndents(m.render(0))
}
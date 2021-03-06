package model

import (
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

type indentedString struct {
	indent int
	s      string
}

func renderIndents(ss []indentedString) string {
	var sb strings.Builder
	for _, l := range ss {
		fmt.Fprint(&sb, strings.Repeat(" ", l.indent*4))
		fmt.Fprintln(&sb, l.s)
	}
	return sb.String()
}

//nolint:deadcode,unused
func renderColors(cs []colorful.Color) string {
	var sb strings.Builder

	sb.WriteString("[")
	for _, c := range cs {
		sb.WriteString(c.Hex())
		sb.WriteString(",")
	}
	sb.WriteString("]")

	return sb.String()
}

func renderLedColors(ls []*LED) string {
	var sb strings.Builder

	sb.WriteString("[")
	for _, l := range ls {
		sb.WriteString(l.newColor.Hex())
		sb.WriteString(",")
	}
	sb.WriteString("]")

	return sb.String()
}

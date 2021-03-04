package model

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

type LED struct {
	zone        *Zone
	mode        Mode // FIXME: really do need to fake that zone up for EffectModes
	name        string
	serverColor colorful.Color
	newColor    colorful.Color
}

func (l *LED) GetName() string {
	// Never seen it not give them unique names, plus they don't have an index on the wire so we'd have to invent it
	return l.name
}

func (l *LED) Size() int {
	return 1
}

func (l *LED) SetColor(c colorful.Color) {
	l.newColor = c
}

func (l *LED) SetColors(cs []colorful.Color) {
	if len(cs) != 1 {
		panic(fmt.Errorf("Trying to set 1 LED with %d colors.", len(cs)))
	}

	l.newColor = cs[0]
}

func (l *LED) Diff() {
	fmt.Printf("%s: %s -> %s\n", l.GetName(), l.serverColor.Hex(), l.newColor.Hex()) // TODO: objects should have * to their parents; this thing should print its whole path
}

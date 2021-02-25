package model

import "github.com/lucasb-eyer/go-colorful"

type LED struct {
	name        string
	serverColor colorful.Color
	newColor    colorful.Color
}

func (l *LED) SetColor(c colorful.Color) {
	l.newColor = c
}

package component

import (
	"image/color"

	"github.com/m110/kingdoms/engine"

	"github.com/yohamta/donburi"
)

type TextSize int

const (
	TextSizeL TextSize = iota
	TextSizeM TextSize = iota
	TextSizeS
)

type TextData struct {
	Text  string
	Color color.Color
	Size  TextSize

	Hidden bool

	Streaming      bool
	StreamingTimer *engine.Timer
}

var Text = donburi.NewComponentType[TextData]()

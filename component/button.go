package component

import "github.com/yohamta/donburi"

type ButtonSound int

const (
	ButtonSoundDefault ButtonSound = iota
	ButtonSoundNone
)

type ButtonData struct {
	// TODO Should have KeepSelection here, not return bool
	OnClick func(w donburi.World, e *donburi.Entry)

	Sound ButtonSound

	// If KeepSelection is true, the currently selected tile won't be unselected
	// once the button is pressed.
	KeepSelection bool
}

var Button = donburi.NewComponentType[ButtonData]()

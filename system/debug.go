package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"

	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
)

type Debug struct {
	debug *component.DebugData
}

func NewDebug() *Debug {
	return &Debug{}
}

func (d *Debug) Init(w donburi.World) {
	d.debug = component.Debug.Get(engine.MustFindWithComponent(w, component.Debug))
}

func (d *Debug) Update(w donburi.World) {
	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		d.debug.Enabled = !d.debug.Enabled
	}

	if !d.debug.Enabled {
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		UpdateResources(w, events.ResourcesDelta{
			Food: 10,
		})
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		UpdateResources(w, events.ResourcesDelta{
			Stone: 10,
		})
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		UpdateResources(w, events.ResourcesDelta{
			Wood: 10,
		})
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		UpdateResources(w, events.ResourcesDelta{
			Iron: 10,
		})
	} else if inpututil.IsKeyJustPressed(ebiten.Key5) {
		UpdateResources(w, events.ResourcesDelta{
			Gold: 1,
		})
	} else if inpututil.IsKeyJustPressed(ebiten.Key6) {
		UpdateResources(w, events.ResourcesDelta{
			Diamonds: 1,
		})
	}
}

package component

import "github.com/yohamta/donburi"

type UIPanelData struct {
	// If true, clicking on this panel will not deselect the tile.
	KeepSelection bool
}

var UIPanel = donburi.NewComponentType[UIPanelData]()

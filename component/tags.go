package component

import "github.com/yohamta/donburi"

var (
	// UI is a UI parent. Doesn't need to be used on children attached to a UI.
	UI                = donburi.NewTag()
	ResearchPanel     = donburi.NewTag()
	SeasonPanel       = donburi.NewTag()
	ResourceInfoPanel = donburi.NewTag()
	UnitInfoPanel     = donburi.NewTag()
	BuildPanel        = donburi.NewTag()

	Road     = donburi.NewTag()
	Settlers = donburi.NewTag()

	SelectedIndicator = donburi.NewTag()

	IconWood     = donburi.NewTag()
	IconStone    = donburi.NewTag()
	IconFood     = donburi.NewTag()
	IconIron     = donburi.NewTag()
	IconGold     = donburi.NewTag()
	IconDiamonds = donburi.NewTag()

	ResourceAmountText = donburi.NewTag()

	TileBorder  = donburi.NewTag()
	TileFog     = donburi.NewTag()
	TileDeposit = donburi.NewTag()

	Dialog         = donburi.NewTag()
	UpkeepDialog   = donburi.NewTag()
	TutorialDialog = donburi.NewTag()
)

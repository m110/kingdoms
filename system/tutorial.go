package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/events"
)

type tutorialStep struct {
	Text           string
	DialogPosition math.Vec2
	OnShow         func(w donburi.World)
	IsComplete     func(w donburi.World, events tutorialEvents) bool
}

var tutorialSteps = []tutorialStep{
	{
		Text: `Welcome to your new kingdom!

You can see your first settlement in the middle.

Look around the map by holding the right mouse
button and moving the mouse.`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			return events.CameraMoved != nil
		},
	},
	{
		Text: `Nice!

Let's build a road.

First, choose the road tool on the right
(or press W).`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			return true
		},
	},
	{
		Text: `Click on the map to build a road.

You need to connect it to one of the roads.`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			return events.RoadSpawned != nil
		},
	},
	{
		Text: `Well done!

Building a road costs 1 stone.

Look out so you don't run out of resources!

Build a road to the mountains tile to collect stone.`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			return events.ResourcesUpdated != nil && events.ResourcesUpdated.Delta.Stone > 0
		},
	},
	{
		Text: `To build a new settlement, you need 10 wood.

Connect a road to the forest tile to collect wood.`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			return events.ResourcesUpdated != nil && events.ResourcesUpdated.Delta.Wood > 0
		},
	},
	{
		Text: `To collect resources, the tile needs
to be in range of a settlement.

You should collect some food, but the berries 
are out of range.

First, build a road to the berries.

Then, select the build tool
(or press E).`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			return true
		},
	},
	{
		Text: `You need to build a settlement on a road
and on a plains tile.
  
Click on the map to build a settlement near berries.`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			return events.SettlementSpawned != nil
		},
	},
	{
		Text: `Great!

Let's upgrade the settlement to a village.

First, choose the select tool
(or press Q).`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			return true
		},
	},
	{
		Text: `Click on one of your settlements.`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			if events.TileSelected == nil {
				return false
			}

			_, ok := transform.FindChildWithComponent(events.TileSelected.Tile, component.Settlement)
			return ok
		},
	},
	{
		Text: `Upgrade costs 1 food.

Click the upgrade button
(or press R).`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			return events.SettlementUpgraded != nil
		},
	},
	{
		Text: `Perfect!

To win the game, you need to build 5 castles
(level 5 settlements).

Tip: when you build a settlement, it gets
a free upgrade for each settlement connected
to it with a straight road!

Now go back to the menu and start the game!

Good luck!`,
		DialogPosition: math.Vec2{
			X: 96,
			Y: 50,
		},
		OnShow: func(w donburi.World) {},
		IsComplete: func(w donburi.World, events tutorialEvents) bool {
			return false
		},
	},
}

type tutorialEvents struct {
	CameraMoved        *events.CameraMoved
	TileSelected       *events.TileSelected
	RoadSpawned        *events.RoadSpawned
	SettlementSpawned  *events.SettlementSpawned
	ResourcesUpdated   *events.ResourcesUpdated
	SettlementUpgraded *events.SettlementUpgraded
}

type Tutorial struct {
	currentStep         int
	tutorialDialogQuery *donburi.Query
}

func NewTutorial() *Tutorial {
	return &Tutorial{
		tutorialDialogQuery: donburi.NewQuery(
			filter.Contains(component.TutorialDialog),
		),
	}
}

func (t *Tutorial) Init(w donburi.World) {
	t.ShowCurrentStep(w)

	events.CameraMovedEvent.Subscribe(w, func(w donburi.World, e events.CameraMoved) {
		t.checkProgress(w, tutorialEvents{
			CameraMoved: &e,
		})
	})
	events.TileSelectedEvent.Subscribe(w, func(w donburi.World, e events.TileSelected) {
		t.checkProgress(w, tutorialEvents{
			TileSelected: &e,
		})
	})
	events.RoadSpawnedEvent.Subscribe(w, func(w donburi.World, e events.RoadSpawned) {
		t.checkProgress(w, tutorialEvents{
			RoadSpawned: &e,
		})
	})
	events.SettlementSpawnedEvent.Subscribe(w, func(w donburi.World, e events.SettlementSpawned) {
		t.checkProgress(w, tutorialEvents{
			SettlementSpawned: &e,
		})
	})
	events.ResourcesUpdatedEvent.Subscribe(w, func(w donburi.World, e events.ResourcesUpdated) {
		t.checkProgress(w, tutorialEvents{
			ResourcesUpdated: &e,
		})
	})
	events.SettlementUpgradedEvent.Subscribe(w, func(w donburi.World, e events.SettlementUpgraded) {
		t.checkProgress(w, tutorialEvents{
			SettlementUpgraded: &e,
		})
	})
}

func (t *Tutorial) Update(w donburi.World) {}

func (t *Tutorial) checkProgress(w donburi.World, events tutorialEvents) {
	if t.currentStep >= len(tutorialSteps) {
		return
	}

	if tutorialSteps[t.currentStep].IsComplete(w, events) {
		t.currentStep++
		t.HideStep(w)
		t.ShowCurrentStep(w)
	}
}

func (t *Tutorial) HideStep(w donburi.World) {
	dialog, ok := t.tutorialDialogQuery.First(w)
	if ok {
		component.Destroy(dialog)
	}
}

func (t *Tutorial) ShowCurrentStep(w donburi.World) {
	if t.currentStep >= len(tutorialSteps) {
		return
	}

	step := tutorialSteps[t.currentStep]

	var buttonEntry *donburi.Entry
	if t.currentStep == len(tutorialSteps)-1 {
		buttonEntry = archetype.NewButton(w, "Close", math.Vec2{
			X: 96,
			Y: 90,
		}, false, func(w donburi.World, e *donburi.Entry) {})
		component.Button.Get(buttonEntry).KeepSelection = true
	}

	dialog := archetype.NewDialogWithButton(w, step.DialogPosition, step.Text, buttonEntry)
	dialog.AddComponent(component.TutorialDialog)
	component.Layer.Get(dialog).Layer = component.SpriteUILayerTutorial

	if buttonEntry != nil {
		button := component.Button.Get(buttonEntry)
		button.KeepSelection = true
		button.OnClick = func(w donburi.World, e *donburi.Entry) {
			component.Destroy(dialog)
		}
	}

	step.OnShow(w)
}

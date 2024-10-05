package system

import (
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
)

type UI struct {
	foodIcon     *donburi.Entry
	stoneIcon    *donburi.Entry
	woodIcon     *donburi.Entry
	ironIcon     *donburi.Entry
	goldIcon     *donburi.Entry
	diamondsIcon *donburi.Entry
}

func NewUI() *UI {
	return &UI{}
}

func (u *UI) Init(w donburi.World) {
	u.foodIcon = engine.MustFindWithComponent(w, component.IconFood)
	u.stoneIcon = engine.MustFindWithComponent(w, component.IconStone)
	u.woodIcon = engine.MustFindWithComponent(w, component.IconWood)
	u.ironIcon = engine.MustFindWithComponent(w, component.IconIron)
	u.goldIcon = engine.MustFindWithComponent(w, component.IconGold)
	u.diamondsIcon, _ = engine.FindWithComponent(w, component.IconDiamonds)

	events.TileSelectedEvent.Subscribe(w, u.onTileSelected)
	events.TileUnselectedEvent.Subscribe(w, u.onTileUnselected)
	events.ResourcesUpdatedEvent.Subscribe(w, u.onResourcesUpdated)
	events.TileResourcesHarvestedEvent.Subscribe(w, u.onTileResourcesHarvested)
	events.SettlementSpawnedEvent.Subscribe(w, u.onSettlementSpawned)
	events.RoadSpawnedEvent.Subscribe(w, u.onRoadSpawned)
	events.OriginSettlementSpawnedEvent.Subscribe(w, u.onOriginSettlementSpawned)
	events.SettlementUpgradedEvent.Subscribe(w, u.onSettlementUpgraded)
	events.GameOverEvent.Subscribe(w, u.onGameOver)
	events.GameWonEvent.Subscribe(w, u.onGameWon)
	events.TechnologyResearchedEvent.Subscribe(w, u.onTechnologyResearched)
	events.ProgressUpdatedEvent.Subscribe(w, u.onProgressUpdated)
}

func (u *UI) Update(w donburi.World) {}

func (u *UI) onTileSelected(w donburi.World, event events.TileSelected) {
	u.unselectTile(w)

	tile := component.Tile.Get(event.Tile)
	if !tile.InSightRange {
		return
	}

	if roadPlaceholdersPresent(w) {
		hideRoadPlaceholders(w)

		if tile.HasRoad {
			showRoadPlaceholders(w, event.Tile)
		}
	}

	_, ok := transform.FindChildWithComponent(event.Tile, component.RoadPlaceholder)
	if ok {
		events.NewRoadRequestedEvent.Publish(w, events.NewRoadRequested{
			Tile: event.Tile,
		})
	}

	// TODO The logic should probably be here?
	archetype.ShowBuildPanelIfRelevant(w, event.Tile)
	archetype.ShowUnitInfoPanelIfRelevant(w, event.Tile)
	archetype.ShowResourceInfoPanelIfRelevant(w, event.Tile)
}

func (u *UI) onTileUnselected(w donburi.World, event events.TileUnselected) {
	u.unselectTile(w)
	hideRoadPlaceholders(w)
}

// TODO what is the place to keep this?
func showRoadPlaceholders(w donburi.World, tileEntry *donburi.Entry) {
	hideRoadPlaceholders(w)

	tile := component.Tile.Get(tileEntry)
	neighborTiles := tile.NeighborTiles.All()

	donburi.NewQuery(filter.Contains(component.RoadPlaceholder)).Each(w, func(entry *donburi.Entry) {
		component.Destroy(entry)
	})

	for i := range neighborTiles {
		n := neighborTiles[i]
		neighborTile := component.Tile.Get(n.Entry)

		if !canBuildRoad(neighborTile) {
			continue
		}

		placeholder := archetype.New(w).
			WithParent(n.Entry).
			WithLayer(component.SpriteLayerRoad).
			WithSprite(component.SpriteData{
				Image: neighborTile.RoadImage(),
				AlphaOverride: &component.AlphaOverride{
					A: 0.5,
				},
			}).
			With(component.RoadPlaceholder).
			With(component.Animation).
			Entry()

		component.RoadPlaceholder.SetValue(placeholder, component.RoadPlaceholderData{
			Direction: n.Direction,
		})

		component.Animation.SetValue(placeholder, component.AnimationData{
			Active: true,
			Timer:  engine.NewTimer(250 * time.Millisecond),
			Update: func(e *donburi.Entry) {
				anim := component.Animation.Get(e)
				if anim.Timer.IsReady() {
					anim.Timer.Reset()
					sprite := component.Sprite.Get(e)
					sprite.Hidden = !sprite.Hidden
				}
			},
		})
	}

	archetype.HideBuildRoadPanel(w)
	if len(neighborTiles) > 0 {
		archetype.ShowBuildRoadPanel(w)
	}
}

func roadPlaceholdersPresent(w donburi.World) bool {
	return donburi.NewQuery(filter.Contains(component.RoadPlaceholder)).Count(w) > 0
}

func hideRoadPlaceholders(w donburi.World) {
	archetype.HideBuildRoadPanel(w)
	donburi.NewQuery(filter.Contains(component.RoadPlaceholder)).Each(w, func(entry *donburi.Entry) {
		component.Destroy(entry)
	})
}

func (u *UI) unselectTile(w donburi.World) {
	archetype.HideUnitInfoPanel(w)
	archetype.HideResourceInfoPanel(w)
	archetype.HideBuildPanel(w)
}

func (u *UI) onResourcesUpdated(w donburi.World, event events.ResourcesUpdated) {
	foodText := engine.MustFindChildWithComponent(u.foodIcon, component.Text)
	component.Text.Get(foodText).Text = strconv.Itoa(event.Resources.Food)

	woodText := engine.MustFindChildWithComponent(u.woodIcon, component.Text)
	component.Text.Get(woodText).Text = strconv.Itoa(event.Resources.Wood)

	stoneText := engine.MustFindChildWithComponent(u.stoneIcon, component.Text)
	component.Text.Get(stoneText).Text = strconv.Itoa(event.Resources.Stone)

	ironText := engine.MustFindChildWithComponent(u.ironIcon, component.Text)
	component.Text.Get(ironText).Text = strconv.Itoa(event.Resources.Iron)

	goldText := engine.MustFindChildWithComponent(u.goldIcon, component.Text)
	component.Text.Get(goldText).Text = strconv.Itoa(event.Resources.Gold)

	if u.diamondsIcon != nil {
		diamondsText := engine.MustFindChildWithComponent(u.diamondsIcon, component.Text)
		component.Text.Get(diamondsText).Text = strconv.Itoa(event.Resources.Diamonds)
	}
}

func (u *UI) onTileResourcesHarvested(w donburi.World, event events.TileResourcesHarvested) {
	var image *ebiten.Image
	if event.Delta.Wood > 0 {
		image = assets.IconWood
	} else if event.Delta.Stone > 0 {
		image = assets.IconStone
	} else if event.Delta.Food > 0 {
		image = assets.IconFood
	} else {
		panic("no resource delta")
	}

	archetype.NewTileFloatingAnimation(w, event.Tile, image)
}

func (u *UI) onSettlementUpgraded(w donburi.World, event events.SettlementUpgraded) {
	tile, ok := transform.GetParent(event.Settlement)
	if !ok {
		panic("parent not found")
	}

	u.updateBuildPanel(w, event.Tile)
	archetype.UpdateUnitInfoPanel(w)
	archetype.NewTileFloatingAnimation(w, tile, assets.IconUpgrade)
}

func (u *UI) onGameOver(w donburi.World, event events.GameOver) {
	archetype.New(w).
		WithPosition(math.Vec2{X: 250, Y: 200}).
		WithLayer(component.SpriteUILayerTop).
		With(component.UI).
		WithText(component.TextData{
			Text: "Game Over",
		})

	archetype.New(w).
		WithPosition(math.Vec2{X: 210, Y: 300}).
		WithLayer(component.SpriteUILayerTop).
		With(component.UI).
		WithText(component.TextData{
			Text: "Out of resources",
		})
}

func (u *UI) onGameWon(w donburi.World, event events.GameWon) {
	archetype.New(w).
		WithPosition(math.Vec2{X: 250, Y: 300}).
		WithLayer(component.SpriteUILayerTop).
		With(component.UI).
		WithText(component.TextData{
			Text: "You did it! Congrats!",
		})
}

func (u *UI) onProgressUpdated(w donburi.World, event events.ProgressUpdated) {
	panel := engine.MustFindWithComponent(w, component.SeasonPanel)
	component.Script.Get(panel).Update(panel)

	archetype.UpdateResourceInfoPanel(w)
}

func (u *UI) onTechnologyResearched(w donburi.World, event events.TechnologyResearched) {
	if event.Technology == domain.TechnologyIronworking {
		component.Active.Get(u.ironIcon).Active = true
	} else if event.Technology == domain.TechnologyCurrency {
		component.Active.Get(u.goldIcon).Active = true
	}

	archetype.UpdateResourceInfoPanel(w)
}

func (u *UI) onRoadSpawned(w donburi.World, event events.RoadSpawned) {
	u.updateBuildPanel(w, event.Tile)
	archetype.UpdateResourceInfoPanel(w)
}

func (u *UI) onSettlementSpawned(w donburi.World, event events.SettlementSpawned) {
	u.updateBuildPanel(w, event.Tile)
	archetype.HideResourceInfoPanel(w)
	archetype.ShowUnitInfoPanelIfRelevant(w, event.Tile)
}

func (u *UI) onOriginSettlementSpawned(w donburi.World, event events.OriginSettlementSpawned) {
	u.updateBuildPanel(w, event.Tile)
	archetype.HideResourceInfoPanel(w)
	archetype.HideUnitInfoPanel(w)
	archetype.ShowUnitInfoPanelIfRelevant(w, event.Tile)
}

func (u *UI) updateBuildPanel(w donburi.World, tileEntry *donburi.Entry) {
	// TODO not sure if could be done better
	if archetype.HideBuildPanel(w) {
		archetype.ShowBuildPanelIfRelevant(w, tileEntry)
	}
}

package system

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
)

const TileImagesTransitionDuration = 1 * time.Second

type SpriteTransition struct {
	query *query.Query
	cache map[*ebiten.Image]map[*ebiten.Image]map[int]*ebiten.Image

	inTransition bool
	timer        *engine.Timer
}

func NewSpriteTransition() *SpriteTransition {
	return &SpriteTransition{
		query: query.NewQuery(filter.Contains(component.SpriteTransition)),
		cache: map[*ebiten.Image]map[*ebiten.Image]map[int]*ebiten.Image{},
		timer: engine.NewTimer(TileImagesTransitionDuration),
	}
}

func (s *SpriteTransition) Init(w donburi.World) {
	events.SeasonChangedEvent.Subscribe(w, s.onSeasonChanged)
}

func (s *SpriteTransition) Update(w donburi.World) {
	if !s.inTransition {
		return
	}

	s.timer.Update()
	if s.timer.IsReady() {
		s.inTransition = false
	}

	s.query.Each(w, func(entry *donburi.Entry) {
		spriteTransition := component.SpriteTransition.Get(entry)
		if !spriteTransition.Active {
			return
		}

		sprite := component.Sprite.Get(entry)
		if s.timer.IsReady() {
			sprite.Image = spriteTransition.To

			spriteTransition.Active = false
			spriteTransition.From = nil
			spriteTransition.To = nil
		} else {
			sprite.Image = s.interpolateImages(spriteTransition.From, spriteTransition.To)
		}
	})

	if !s.inTransition {
		clear(s.cache)
	}
}

func (s *SpriteTransition) onSeasonChanged(w donburi.World, event events.SeasonChanged) {
	s.timer.Reset()
	s.inTransition = true

	query.NewQuery(filter.Contains(component.Tile)).Each(w, func(entry *donburi.Entry) {
		tile := component.Tile.Get(entry)
		terrainData := domain.Terrains[tile.Terrain]

		tileSprite := component.Sprite.Get(entry)
		tileSpriteTransition := component.SpriteTransition.Get(entry)
		newTerrainImage := terrainData.Sprites.BySeason(event.Progress.Season)

		if tile.InSightRange {
			tileSpriteTransition.Active = true
			tileSpriteTransition.From = tileSprite.Image
			tileSpriteTransition.To = newTerrainImage
		} else {
			tileSpriteTransition.Active = false
			tileSprite.Image = newTerrainImage
		}

		if tile.Deposit == nil {
			return
		}

		updateDepositImage(w, entry)
	})
}

func (s *SpriteTransition) onTileResourcesHarvested(w donburi.World, event events.TileResourcesHarvested) {
	// TODO timer doesn't start because it's global
	updateDepositImage(w, event.Tile)
}

func updateDepositImage(w donburi.World, tileEntry *donburi.Entry) {
	progress := engine.MustFindComponent[component.ProgressData](w, component.Progress)
	tile := component.Tile.Get(tileEntry)

	deposit := engine.MustFindChildWithComponent(tileEntry, component.TileDeposit)
	depositSprite := component.Sprite.Get(deposit)
	depositSpriteTransition := component.SpriteTransition.Get(deposit)
	newDepositImage := tile.DepositImage(progress.Season)

	if tile.InSightRange {
		depositSpriteTransition.Active = true
		depositSpriteTransition.From = depositSprite.Image
		depositSpriteTransition.To = newDepositImage
	} else {
		depositSpriteTransition.Active = false
		depositSprite.Image = newDepositImage
	}
}

func (s *SpriteTransition) interpolateImages(from *ebiten.Image, to *ebiten.Image) *ebiten.Image {
	if from.Bounds().Size() != to.Bounds().Size() {
		panic("images must have the same size")
	}

	if _, ok := s.cache[from]; !ok {
		s.cache[from] = map[*ebiten.Image]map[int]*ebiten.Image{}
	}

	if _, ok := s.cache[from][to]; !ok {
		s.cache[from][to] = map[int]*ebiten.Image{}
	}

	factor := s.timer.PercentDone()
	intFactor := int(factor * 100)

	targetSprite, ok := s.cache[from][to][intFactor]
	if ok {
		return targetSprite
	}

	targetSprite = ebiten.NewImageFromImage(from)
	s.cache[from][to][intFactor] = targetSprite

	for y := 0; y < from.Bounds().Dy(); y++ {
		for x := 0; x < from.Bounds().Dx(); x++ {
			r1, g1, b1, a1 := normalizedColor(from.At(x, y))
			r2, g2, b2, a2 := normalizedColor(to.At(x, y))

			if a1 == 0 && a2 == 0 {
				continue
			}

			if r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2 {
				continue
			}

			r := uint8(float64(r1) + (float64(r2)-float64(r1))*factor)
			g := uint8(float64(g1) + (float64(g2)-float64(g1))*factor)
			b := uint8(float64(b1) + (float64(b2)-float64(b1))*factor)
			a := uint8(float64(a1) + (float64(a2)-float64(a1))*factor)

			targetSprite.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	return targetSprite
}

func normalizedColor(c color.Color) (r, g, b, a uint32) {
	r, g, b, a = c.RGBA()
	r = r >> 8
	g = g >> 8
	b = b >> 8
	a = a >> 8
	return r, g, b, a
}

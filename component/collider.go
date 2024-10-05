package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/kingdoms/engine"
)

const (
	CollisionLayerPlayers ColliderLayer = iota
	CollisionLayerTiles
	CollisionLayerUI
	CollisionLayerButtons
)

type ColliderLayer int

type ColliderData struct {
	Offset math.Vec2

	Width  float64
	Height float64
	Layer  ColliderLayer

	Collisions             map[CollisionKey]Collision
	JustCollidedWith       map[CollisionKey]struct{}
	JustOutOfCollisionWith map[CollisionKey]struct{}
}

func (d *ColliderData) Rect(entry *donburi.Entry) engine.Rect {
	pos := transform.WorldPosition(entry)
	scale := transform.WorldScale(entry)
	return engine.Rect{
		X:      pos.X + d.Offset.X,
		Y:      pos.Y + d.Offset.Y,
		Width:  d.Width * scale.X,
		Height: d.Height * scale.Y,
	}
}

type CollisionKey struct {
	Layer ColliderLayer
	Other donburi.Entity
}

type Collision struct {
	TimesSeen int
	Detected  bool
}

var Collider = donburi.NewComponentType[ColliderData]()

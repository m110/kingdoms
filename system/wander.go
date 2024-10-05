package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/kingdoms/component"
)

type Wander struct {
	query *query.Query
}

func NewWander() *Wander {
	return &Wander{
		query: query.NewQuery(
			filter.Contains(
				component.Wander,
				component.Velocity,
			),
		),
	}
}

func (w *Wander) Init(world donburi.World) {}

func (w *Wander) Update(world donburi.World) {
	w.query.Each(world, func(entry *donburi.Entry) {
		pos := transform.GetTransform(entry).LocalPosition
		wander := component.Wander.Get(entry)
		velocity := component.Velocity.Get(entry)

		// TODO naive implementation but good enough for now
		if pos.X > wander.Max.X || pos.X < wander.Min.X || pos.Y > wander.Max.Y || pos.Y < wander.Min.Y {
			velocity.Velocity.X *= -1
			velocity.Velocity.Y *= -1
		}
	})
}

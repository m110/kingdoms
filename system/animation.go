package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/kingdoms/component"
)

type Animation struct {
	query *query.Query
}

func NewAnimation() *Animation {
	return &Animation{
		query: query.NewQuery(filter.Contains(component.Animation)),
	}
}

func (s *Animation) Init(w donburi.World) {}

func (s *Animation) Update(w donburi.World) {
	s.query.Each(w, func(entry *donburi.Entry) {
		animation := component.Animation.Get(entry)
		if !animation.Active {
			return
		}
		if animation.Timer != nil {
			animation.Timer.Update()
		}
		animation.Update(entry)
	})
}

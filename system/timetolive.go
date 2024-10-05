package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/kingdoms/component"
)

type TimeToLive struct {
	query *query.Query
}

func NewTimeToLive() *TimeToLive {
	return &TimeToLive{
		query: query.NewQuery(
			filter.Contains(component.TimeToLive),
		),
	}
}

func (t *TimeToLive) Init(w donburi.World) {}

func (t *TimeToLive) Update(w donburi.World) {
	t.query.Each(w, func(entry *donburi.Entry) {
		ttl := component.TimeToLive.Get(entry)
		ttl.Timer.Update()
		if ttl.Timer.IsReady() {
			component.Destroy(entry)
		}
	})
}

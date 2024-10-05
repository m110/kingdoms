package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/kingdoms/component"
)

type Text struct {
	query *query.Query
}

func NewText() *Text {
	return &Text{
		query: query.NewQuery(
			filter.Contains(
				transform.Transform,
				component.Text,
			),
		),
	}
}

func (t *Text) Init(w donburi.World) {}

func (t *Text) Update(w donburi.World) {
	t.query.Each(w, func(entry *donburi.Entry) {
		txt := component.Text.Get(entry)
		if txt.Streaming && !txt.StreamingTimer.IsReady() {
			txt.StreamingTimer.Update()
		}
	})
}

package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"

	"github.com/m110/kingdoms/assets"

	"github.com/m110/kingdoms/component"
)

func NewTable(
	w donburi.World,
	rows int,
	columns int,
	cellWidth int,
	cellHeight int,
	headerWidth int,
	headerHeight int,
) EntryBuilder {
	borderSize := 3
	totalWidth := headerWidth + cellWidth*(columns-1) + borderSize*(columns-1)
	totalHeight := headerHeight + cellHeight*(rows-1) + borderSize*(rows-1)
	backgroundImage := ebiten.NewImage(totalWidth, totalHeight)

	for i := 1; i < rows; i++ {
		vector.DrawFilledRect(
			backgroundImage,
			0,
			float32(headerHeight+cellHeight*(i-1)+borderSize*(i-1)),
			float32(totalWidth),
			float32(borderSize),
			assets.TableBorderColor,
			true,
		)
	}

	for i := 1; i < columns; i++ {
		vector.DrawFilledRect(
			backgroundImage,
			float32(headerWidth+cellWidth*(i-1)+borderSize*(i-1)),
			0,
			float32(borderSize),
			float32(totalHeight),
			assets.TableBorderColor,
			true,
		)
	}

	table := New(w).
		WithSprite(component.SpriteData{
			Image: backgroundImage,
		}).
		With(component.Table)

	tableData := component.Table.Get(table.entry)
	tableData.Rows = make([][]*donburi.Entry, rows)
	for i := range tableData.Rows {
		tableData.Rows[i] = make([]*donburi.Entry, columns)
		for j := range tableData.Rows[i] {
			var x, y int
			if j > 0 {
				x = cellWidth*(j-1) + headerWidth
			}
			if i > 0 {
				y = cellHeight*(i-1) + headerHeight
			}

			tableData.Rows[i][j] = New(w).
				WithParent(table.Entry()).
				WithPosition(math.Vec2{
					X: float64(x) + float64(borderSize*j),
					Y: float64(y) + float64(borderSize*i),
				}).
				Entry()
		}
	}

	return table
}

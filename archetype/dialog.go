package archetype

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"golang.org/x/image/font"

	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
)

const (
	frameBorderWidth = 3
)

func NewDialog(
	w donburi.World,
	pos math.Vec2,
	content string,
) *donburi.Entry {
	return NewDialogWithButton(w, pos, content, nil)
}

func NewDialogWithButton(
	w donburi.World,
	pos math.Vec2,
	content string,
	button *donburi.Entry,
) *donburi.Entry {
	marginTop := 20
	marginBottom := 10
	marginLeft := 10
	lineHeight := 20

	if button != nil {
		marginBottom += 60
	}

	face := assets.MediumSquareFont

	maxWidth := 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		// TODO 60 picked arbitrarily here, not sure why MeasureString returns huge values
		width := int(font.MeasureString(face, line)) / 60
		if width > maxWidth {
			maxWidth = width
		}
	}

	totalWidth := maxWidth + marginLeft
	totalHeight := (len(lines)-1)*lineHeight + marginTop + marginBottom

	frame := NewFrameWithBorder(w, pos, totalWidth, totalHeight)
	frame.AddComponent(component.Dialog)

	image := component.Sprite.Get(frame).Image

	for i, line := range lines {
		text.Draw(
			image,
			line,
			face,
			marginLeft+frameBorderWidth,
			marginTop+frameBorderWidth+lineHeight*i,
			assets.TextColor,
		)
	}

	if button != nil {
		buttonTransform := transform.GetTransform(button)
		// TODO naive calculations
		buttonTransform.LocalPosition = math.Vec2{
			X: float64(totalWidth/2) - 20,
			Y: float64(totalHeight - 40),
		}
		transform.AppendChild(frame, button, false)
	}

	return frame
}

func NewFrame(w donburi.World, pos math.Vec2, width int, height int) *donburi.Entry {
	image := ebiten.NewImage(width, height)
	vector.DrawFilledRect(image, 0, 0, float32(width), float32(height), assets.PanelColor, true)

	return newFrame(w, pos, image, float64(width), float64(height))
}

func NewFrameWithBorder(w donburi.World, pos math.Vec2, width int, height int) *donburi.Entry {
	totalWidth := width + frameBorderWidth*2
	totalHeight := height + frameBorderWidth*2

	image := ebiten.NewImage(totalWidth, totalHeight)
	image.Fill(assets.BorderColor)
	vector.DrawFilledRect(
		image,
		float32(frameBorderWidth), float32(frameBorderWidth),
		float32(totalWidth-frameBorderWidth*2), float32(totalHeight-frameBorderWidth*2),
		assets.PanelColor,
		true,
	)

	return newFrame(w, pos, image, float64(totalWidth), float64(totalHeight))
}

func newFrame(w donburi.World, pos math.Vec2, image *ebiten.Image, width float64, height float64) *donburi.Entry {
	frame := New(w).
		WithPosition(pos).
		WithLayer(component.SpriteUILayerDialog).
		WithSprite(component.SpriteData{
			Image: image,
		}).
		With(component.Collider).
		With(component.UI).
		Entry()

	component.Collider.SetValue(frame, component.ColliderData{
		Width:  width,
		Height: height,
	})

	return frame
}

package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/m110/kingdoms/domain"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
)

const (
	researchButtonSize = 48
)

func NewResearchPanel(
	w donburi.World,
	selectedTech *domain.Technology,
) *donburi.Entry {
	panel := NewFrameWithBorder(w, math.Vec2{X: 116, Y: 48}, 577, 518)
	panel.AddComponent(component.ResearchPanel)
	panel.AddComponent(component.UIPanel)

	closeButton := NewButton(w, "X", math.Vec2{X: 525, Y: 10}, false, func(w donburi.World, e *donburi.Entry) {
		component.Destroy(panel)
	})
	component.Button.Get(closeButton).KeepSelection = true
	transform.AppendChild(panel, closeButton, false)

	New(w).
		WithPosition(math.Vec2{X: 230, Y: 40}).
		WithParent(panel).
		WithText(component.TextData{
			Text: "Research",
		})

	researchSection := newStartResearchSection(w)
	transform.AppendChild(panel, researchSection.entry, false)

	for i, tech := range domain.RootTechnologies {
		addResearchButton(
			w,
			panel,
			tech,
			80.0,
			180.0,
			i,
			0,
			researchSection,
		)
	}

	if selectedTech != nil {
		researchSection.Select(w, domain.Technologies[*selectedTech])
	}

	return panel
}

type startResearchSection struct {
	selectedTechnology *domain.Technology
	entry              *donburi.Entry
	titleText          *component.TextData
	descriptionText    *component.TextData
	researchButton     *donburi.Entry
	costSection        *donburi.Entry
}

func newStartResearchSection(w donburi.World) *startResearchSection {
	section := New(w).
		WithPosition(math.Vec2{X: 30, Y: 400}).
		With(component.Active).
		Entry()

	component.Active.SetValue(section, component.ActiveData{
		Active: false,
	})

	titleText := New(w).
		WithParent(section).
		WithText(component.TextData{
			Text: "",
		}).Entry()

	descriptionText := New(w).
		WithParent(section).
		WithPosition(math.Vec2{X: 0, Y: 20}).
		WithText(component.TextData{
			Text:  "",
			Size:  component.TextSizeM,
			Color: assets.TextSecondaryColor,
		}).Entry()

	costSection := New(w).
		WithParent(section).
		WithPosition(math.Vec2{X: 0, Y: 50}).
		Entry()

	s := &startResearchSection{
		entry:           section,
		titleText:       component.Text.Get(titleText),
		descriptionText: component.Text.Get(descriptionText),
		costSection:     costSection,
	}

	researchButton := NewButton(w, "Research", math.Vec2{X: 360, Y: 50}, true, func(w donburi.World, e *donburi.Entry) {
		if s.selectedTechnology == nil {
			return
		}

		events.ResearchRequestedEvent.Publish(w, events.ResearchRequested{
			Technology: *s.selectedTechnology,
		})
	})
	transform.AppendChild(section, researchButton, false)

	s.researchButton = researchButton

	return s
}

func (s *startResearchSection) Select(w donburi.World, tech domain.TechnologyData) {
	s.selectedTechnology = &tech.ID
	s.titleText.Text = tech.Name
	s.descriptionText.Text = tech.Description
	component.Active.SetValue(s.entry, component.ActiveData{
		Active: true,
	})

	children, ok := transform.GetChildren(s.costSection)
	if ok {
		for _, child := range children {
			component.Destroy(child)
		}
	}

	research := component.Research.Get(engine.MustFindWithComponent(w, component.Research))
	techKnown := research.KnownTechnologies[*s.selectedTechnology]

	if !techKnown {
		marginX := 80.0
		i := 0

		type cost struct {
			resource domain.Resource
			value    int
		}

		costs := []cost{
			{resource: domain.ResourceFood, value: tech.Cost.Food},
			{resource: domain.ResourceStone, value: tech.Cost.Stone},
			{resource: domain.ResourceWood, value: tech.Cost.Wood},
			{resource: domain.ResourceIron, value: tech.Cost.Iron},
			{resource: domain.ResourceGold, value: tech.Cost.Gold},
		}

		for _, c := range costs {
			if c.value <= 0 {
				continue
			}

			icon := NewResourceIconWithValue(
				w,
				c.resource,
				c.value,
				math.Vec2{X: float64(i) * marginX, Y: 0},
			)
			transform.AppendChild(s.costSection, icon, false)

			i++
		}
	}

	canResearch := !techKnown
	// Check if all technologies that enable this technology are known
	if canResearch {
		for _, tech := range domain.Technologies {
			for _, e := range tech.Enables {
				if e == *s.selectedTechnology && !research.KnownTechnologies[tech.ID] {
					canResearch = false
					break
				}
			}
		}
	}

	if canResearch {
		setButtonActive(w, s.researchButton, true)
	} else {
		component.Active.Get(s.researchButton).Active = false
	}
}

// TODO refactor into factory or something, this is terrible
func addResearchButton(
	w donburi.World,
	parent *donburi.Entry,
	tech domain.Technology,
	startX float64,
	marginX float64,
	xPos int,
	yPos int,
	section *startResearchSection,
) {
	baseY := 60.0
	marginY := 85.0

	x := startX + float64(xPos)*marginX
	y := baseY + float64(yPos)*marginY

	pos := math.Vec2{
		X: x,
		Y: y,
	}
	techData := domain.Technologies[tech]

	button := newResearchButton(w, pos, techData, section)
	transform.AppendChild(parent, button, false)

	research := component.Research.Get(engine.MustFindWithComponent(w, component.Research))

	frameColor := assets.IconFrameColor
	if research.KnownTechnologies[tech] {
		frameColor = assets.IconFrameLightColor
	}

	frameWidth := 2.0
	parentImage := component.Sprite.Get(parent).Image
	vector.StrokeRect(
		parentImage,
		float32(x-frameWidth),
		float32(y-frameWidth),
		float32(researchButtonSize+frameWidth*2),
		float32(researchButtonSize+frameWidth*2),
		float32(frameWidth*2),
		frameColor,
		true,
	)

	lineWidth := 4.0

	// TODO implement for 3 or more if needed
	// (simple for now, whatever)
	if len(techData.Enables) == 1 {
		addResearchButton(w, parent, techData.Enables[0], x, marginX, 0, yPos+1, section)

		lineX := x + researchButtonSize/2
		vector.StrokeLine(
			parentImage,
			float32(lineX),
			float32(y+researchButtonSize),
			float32(lineX),
			float32(y+marginY-lineWidth),
			float32(lineWidth),
			frameColor,
			true,
		)
	} else if len(techData.Enables) == 2 {
		addResearchButton(w, parent, techData.Enables[0], x-marginX/4, marginX/2, 0, yPos+1, section)
		addResearchButton(w, parent, techData.Enables[1], x-marginX/4, marginX/2, 1, yPos+1, section)

		lineX := x + researchButtonSize/2
		lineY := y + researchButtonSize + marginY/4

		leftLineX := x - marginX/4 + researchButtonSize/2
		rightLineX := x + marginX/4 + researchButtonSize/2

		vector.StrokeLine(
			parentImage,
			float32(lineX),
			float32(y+researchButtonSize),
			float32(lineX),
			float32(lineY),
			float32(lineWidth),
			frameColor,
			true,
		)

		vector.StrokeLine(
			parentImage,
			float32(leftLineX),
			float32(lineY),
			float32(rightLineX),
			float32(lineY),
			float32(lineWidth),
			frameColor,
			true,
		)

		vector.StrokeLine(
			parentImage,
			float32(leftLineX),
			float32(lineY-lineWidth/2),
			float32(leftLineX),
			float32(y+marginY-lineWidth),
			float32(lineWidth),
			frameColor,
			true,
		)

		vector.StrokeLine(
			parentImage,
			float32(rightLineX),
			float32(lineY-lineWidth/2),
			float32(rightLineX),
			float32(y+marginY-lineWidth),
			float32(lineWidth),
			frameColor,
			true,
		)
	}
}

func newResearchButton(
	w donburi.World,
	pos math.Vec2,
	technology domain.TechnologyData,
	section *startResearchSection,
) *donburi.Entry {
	image := ebiten.NewImage(researchButtonSize, researchButtonSize)
	image.Fill(assets.ButtonColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(3, 3)
	image.DrawImage(technology.Icon, op)
	image.DrawImage(assets.IconFrame, op)

	button := New(w).
		WithPosition(pos).
		WithSprite(component.SpriteData{
			Image: image,
		}).
		With(component.Collider).
		With(component.Button).
		Entry()

	component.Collider.SetValue(button, component.ColliderData{
		Width:  float64(researchButtonSize),
		Height: float64(researchButtonSize),
		Layer:  component.CollisionLayerButtons,
	})
	component.Button.SetValue(button, component.ButtonData{
		OnClick: func(w donburi.World, e *donburi.Entry) {
			section.Select(w, technology)
		},
	})

	return button
}

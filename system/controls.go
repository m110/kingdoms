package system

import (
	"fmt"
	"image/color"
	stdmath "math"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/colornames"

	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
)

const (
	zoomStep      = 0.05
	touchZoomStep = 0.03

	scrollTriggerDist = 10

	touchDuration = 200 * time.Millisecond
)

type Controls struct {
	mapInteractive bool

	scrollCursorStart math.Vec2
	scrollTouchID     ebiten.TouchID
	mapScrolled       bool
	startedScrolling  bool
	stoppedScrolling  bool

	zoomFactor float64

	cameraBoundsMin math.Vec2
	cameraBoundsMax math.Vec2

	chunksQuery    *query.Query
	crosshairQuery *query.Query
	selectedQuery  *query.Query
	buttonsQuery   *query.Query
	uiPanelsQuery  *query.Query

	isTouching       bool
	isTouchingUI     bool
	isTouchScrolling bool
	touchPos         math.Vec2
	touchTimer       *engine.Timer

	touchID1      ebiten.TouchID
	touchID2      ebiten.TouchID
	pinchDistance float64
	isPinching    bool

	board *component.BoardData
	game  *component.GameData
	debug *component.DebugData
}

func NewControls(mapInteractive bool) *Controls {
	return &Controls{
		mapInteractive: mapInteractive,
		cameraBoundsMin: math.Vec2{
			X: 0,
			Y: 0,
		},
		chunksQuery: query.NewQuery(
			filter.Contains(
				transform.Transform,
				component.Chunk,
			),
		),
		crosshairQuery: query.NewQuery(filter.Contains(component.Crosshair)),
		selectedQuery:  query.NewQuery(filter.Contains(component.SelectedIndicator)),
		buttonsQuery: query.NewQuery(
			filter.Contains(
				component.Collider,
				component.Button,
			),
		),
		uiPanelsQuery: query.NewQuery(
			filter.Or(
				filter.Contains(
					component.UIPanel,
				),
				filter.Contains(
					component.Dialog,
				),
			),
		),
		touchTimer: engine.NewTimer(touchDuration),
	}
}

func (c *Controls) Init(w donburi.World) {
	game := component.MustFindGame(w)
	board := archetype.MustFindBoard(w)

	mapWidth := board.Size.Width * board.TileSize.Width
	mapHeight := board.Size.Height * board.TileSize.Height

	screenWidth := game.Settings.ScreenWidth
	screenHeight := game.Settings.ScreenHeight

	// TODO should be simply based on Board's image?
	c.cameraBoundsMax = math.Vec2{
		X: float64(mapWidth - screenWidth),
		Y: float64(mapHeight - screenHeight),
	}

	c.board = board
	c.game = game
	c.debug = component.Debug.Get(engine.MustFindWithComponent(w, component.Debug))

	camera := archetype.MustFindCamera(w)
	cam := component.Camera.Get(camera)
	transform.GetTransform(camera).LocalScale = math.Vec2{
		X: cam.Zoom.Max,
		Y: cam.Zoom.Max,
	}
}

func (c *Controls) Update(w donburi.World) {
	var cursorX, cursorY int
	var clicked bool
	var uiPanelUnderCursor *donburi.Entry

	if UseTouchControls() {
		touchIDs := ebiten.AppendTouchIDs(nil)

		c.zoomFactor = 0

		switch len(touchIDs) {
		case 0:
			if c.isTouching {
				if c.isTouchingUI {
					c.isTouching = false
					c.isTouchingUI = false
				} else {
					// Finger lifted within given duration, single tap
					if !c.touchTimer.IsReady() {
						c.isTouching = false
						clicked = true
						cursorX, cursorY = int(c.touchPos.X), int(c.touchPos.Y)
					}
				}
			}

			if c.isTouchScrolling {
				c.isTouchScrolling = false
				c.startedScrolling = false
				c.stoppedScrolling = true
			}

			if c.isPinching {
				c.isPinching = false
			}
		case 1:
			if c.isTouching {
				// Start scrolling
				if !c.isTouchingUI {
					cx, cy := ebiten.TouchPosition(touchIDs[0])
					cursor := math.Vec2{X: float64(cx), Y: float64(cy)}
					dist := c.touchPos.Distance(cursor)
					c.touchTimer.Update()
					if c.touchTimer.IsReady() || dist > scrollTriggerDist {
						c.isTouching = false
						c.scrollTouchID = touchIDs[0]
						c.isTouchScrolling = true
						c.startedScrolling = true
						c.stoppedScrolling = false
						cursorX, cursorY = ebiten.TouchPosition(touchIDs[0])
					}
				}
			} else {
				if c.isTouchScrolling {
					c.startedScrolling = false
					c.stoppedScrolling = false
					cursorX, cursorY = ebiten.TouchPosition(c.scrollTouchID)
				} else {
					cursorX, cursorY = ebiten.TouchPosition(touchIDs[0])
					uiPanelUnderCursor = c.getUIPanelUnderCursor(w, cursorX, cursorY)

					if c.debug.Enabled {
						spawnClickDebug(w, cursorX, cursorY, colornames.Limegreen)
					}

					c.isTouching = true
					c.touchPos = math.Vec2{X: float64(cursorX), Y: float64(cursorY)}

					if uiPanelUnderCursor == nil {
						// Map touched
						c.touchTimer.Reset()
					} else {
						// UI Panel touched
						clicked = true
						cursorX, cursorY = int(c.touchPos.X), int(c.touchPos.Y)
						c.isTouchingUI = true
					}
				}
			}

			if c.isPinching {
				c.isPinching = false
			}
		case 2:
			if c.isTouching {
				c.isTouching = false
				c.isTouchingUI = false
			}

			if c.isTouchScrolling {
				c.isTouchScrolling = false
				c.stoppedScrolling = true
			}

			if c.isPinching {
				dist := distanceBetweenTouchIDs(touchIDs[0], touchIDs[1])
				if dist != 0 {
					if dist > c.pinchDistance {
						c.zoomFactor = dist / c.pinchDistance * touchZoomStep
					} else if dist < c.pinchDistance {
						c.zoomFactor = c.pinchDistance / dist * -touchZoomStep
					}
					c.pinchDistance = dist
				}
			} else if c.mapInteractive {
				cx1, cy1 := ebiten.TouchPosition(touchIDs[0])
				cx2, cy2 := ebiten.TouchPosition(touchIDs[1])

				panel1 := c.getUIPanelUnderCursor(w, cx1, cy1)
				panel2 := c.getUIPanelUnderCursor(w, cx2, cy2)

				if panel1 == nil && panel2 == nil {
					c.isPinching = true
					c.pinchDistance = distanceBetweenTouchIDs(touchIDs[0], touchIDs[1])
				}
			}
		default:
			// Multi-touch gesture
			if c.isTouching {
				c.isTouching = false
				c.isTouchingUI = false
			}

			if c.isTouchScrolling {
				c.isTouchScrolling = false
				c.stoppedScrolling = true
			}

			if c.isPinching {
				c.isPinching = false
			}
		}
	} else {
		cursorX, cursorY = ebiten.CursorPosition()
		clicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

		uiPanelUnderCursor = c.getUIPanelUnderCursor(w, cursorX, cursorY)

		if c.mapInteractive {
			c.startedScrolling = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
			c.stoppedScrolling = inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight)

			_, wheelY := ebiten.Wheel()
			if wheelY > 0 {
				c.zoomFactor = zoomStep
			} else if wheelY < 0 {
				c.zoomFactor = -zoomStep
			} else {
				c.zoomFactor = 0
			}
		}
	}

	if c.debug.Enabled && clicked {
		spawnClickDebug(w, cursorX, cursorY, colornames.Magenta)
	}

	cursor := math.Vec2{X: float64(cursorX), Y: float64(cursorY)}

	camera := archetype.MustFindCamera(w)
	cameraPos := transform.GetTransform(camera).LocalPosition
	cameraScale := transform.GetTransform(camera).LocalScale

	c.updateScrolling(w, uiPanelUnderCursor != nil)
	if c.mapScrolled {
		return
	}

	previousCrosshair, _ := c.crosshairQuery.First(w)
	keepCrosshair := false

	if uiPanelUnderCursor != nil {
		if clicked {
			// Handle buttons
			var found bool

			buttons := engine.FindChildrenWithComponent(uiPanelUnderCursor, component.Button)

			for _, entry := range buttons {
				if entry.HasComponent(component.Active) {
					active := component.Active.Get(entry)
					if !active.Active {
						continue
					}
				}

				collider := component.Collider.Get(entry)
				rect := collider.Rect(entry)
				if rect.Contains(cursor) {
					button := component.Button.Get(entry)
					button.OnClick(w, entry)
					if !button.KeepSelection {
						c.unselect(w)
					}

					found = true
					break
				}
			}

			if !found {
				keepSelection := false
				if uiPanelUnderCursor.HasComponent(component.UIPanel) {
					uiPanel := component.UIPanel.Get(uiPanelUnderCursor)
					keepSelection = uiPanel.KeepSelection
				}

				if !keepSelection {
					c.unselect(w)
				}
			}
		}
	} else if c.mapInteractive {
		found := false

		// TODO deduplicate with render
		cameraWidth := float64(c.game.Settings.ScreenWidth) / cameraScale.X
		cameraHeight := float64(c.game.Settings.ScreenHeight) / cameraScale.Y

		cameraEdgeLeft := cameraPos.X
		cameraEdgeRight := cameraPos.X + cameraWidth
		cameraEdgeTop := cameraPos.Y
		cameraEdgeBottom := cameraPos.Y + cameraHeight

		cameraRect := engine.NewRect(cameraEdgeLeft, cameraEdgeTop, cameraWidth, cameraHeight)

		c.chunksQuery.Each(w, func(entry *donburi.Entry) {
			if found {
				return
			}

			chunk := component.Chunk.Get(entry)
			pos := transform.WorldPosition(entry)
			rect := engine.NewRect(pos.X, pos.Y, float64(chunk.Size.Width), float64(chunk.Size.Height))

			if !rect.Intersects(cameraRect) {
				return
			}

			tiles := engine.FindChildrenWithComponent(entry, component.Tile)
			for _, tileEntry := range tiles {
				position := transform.WorldPosition(tileEntry)

				if position.X < cameraEdgeLeft-float64(c.board.TileSize.Width) ||
					position.X > cameraEdgeRight ||
					position.Y < cameraEdgeTop-float64(c.board.TileSize.Height) ||
					position.Y > cameraEdgeBottom {
					continue
				}

				collider := component.Collider.Get(tileEntry)

				tileRect := collider.Rect(tileEntry)

				tileRect.X -= cameraPos.X
				tileRect.Y -= cameraPos.Y

				tileRect.X *= cameraScale.X
				tileRect.Y *= cameraScale.Y
				tileRect.Width *= cameraScale.X
				tileRect.Height *= cameraScale.Y

				if tileRect.Contains(cursor) {
					found = true

					tile := component.Tile.Get(tileEntry)
					if previousCrosshair != nil {
						crosshair := component.Crosshair.Get(previousCrosshair)
						if tile.Position.X == crosshair.Position.X && tile.Position.Y == crosshair.Position.Y {
							keepCrosshair = true
						}
					}

					if !UseTouchControls() && !keepCrosshair {
						c.showCrosshair(w, tileEntry, tile.Position)
					}

					if clicked {
						c.selectTile(w, tileEntry)
					}

					return
				}
			}
		})

	}

	if previousCrosshair != nil && !keepCrosshair {
		component.Destroy(previousCrosshair)
	}

	if uiPanelUnderCursor == nil && c.zoomFactor != 0 {
		cameraWidth := float64(c.game.Settings.ScreenWidth) / cameraScale.X
		cameraHeight := float64(c.game.Settings.ScreenHeight) / cameraScale.Y

		centerX := cameraPos.X + cameraWidth/2.0
		centerY := cameraPos.Y + cameraHeight/2.0

		cameraScale.X += c.zoomFactor
		cameraScale.Y += c.zoomFactor

		cam := component.Camera.Get(camera)

		cameraScale = engine.ClampVector(
			cameraScale,
			math.Vec2{X: cam.Zoom.Min, Y: cam.Zoom.Min},
			math.Vec2{X: cam.Zoom.Max, Y: cam.Zoom.Max},
		)
		transform.GetTransform(camera).LocalScale = cameraScale

		cameraWidth = float64(c.game.Settings.ScreenWidth) / cameraScale.X
		cameraHeight = float64(c.game.Settings.ScreenHeight) / cameraScale.Y

		cameraPos.X = centerX - cameraWidth/2.0
		cameraPos.Y = centerY - cameraHeight/2.0
		transform.GetTransform(camera).LocalPosition = cameraPos
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		c.unselect(w)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		selected, ok := c.selectedQuery.First(w)
		if ok {
			parent, ok := transform.GetParent(selected)
			if !ok {
				panic("selected indicator has no parent")
			}

			sett, ok := transform.FindChildWithComponent(parent, component.Settlement)
			if ok {
				events.SettlementUpgradeRequestedEvent.Publish(w, events.SettlementUpgradeRequested{
					Settlement: sett,
				})
			}
		}
	}
}

func (c *Controls) getUIPanelUnderCursor(w donburi.World, x int, y int) *donburi.Entry {
	uiPanelsByLayer := make(map[int][]*donburi.Entry)
	c.uiPanelsQuery.Each(w, func(entry *donburi.Entry) {
		l := int(component.Layer.Get(entry).Layer)
		uiPanelsByLayer[l] = append(uiPanelsByLayer[l], entry)
	})

	cursor := math.Vec2{X: float64(x), Y: float64(y)}

	var layers []int
	for layer := range uiPanelsByLayer {
		layers = append(layers, int(layer))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(layers)))

	for _, i := range layers {
		for _, entry := range uiPanelsByLayer[i] {
			if entry.HasComponent(component.Active) {
				active := component.Active.Get(entry)
				if !active.Active {
					continue
				}
			}

			collider := component.Collider.Get(entry)
			rect := collider.Rect(entry)
			if rect.Contains(cursor) {
				return entry
			}
		}
	}

	return nil
}

func (c *Controls) updateScrolling(w donburi.World, cursorOverUI bool) {
	if !c.mapInteractive {
		return
	}

	var cursorX, cursorY int

	if UseTouchControls() {
		cursorX, cursorY = ebiten.TouchPosition(c.scrollTouchID)
	} else {
		cursorX, cursorY = ebiten.CursorPosition()
	}

	cursor := math.Vec2{X: float64(cursorX), Y: float64(cursorY)}

	if c.startedScrolling {
		if !cursorOverUI {
			c.mapScrolled = true
			c.scrollCursorStart = cursor
		}
	}

	if c.stoppedScrolling {
		c.mapScrolled = false
	}

	if c.mapScrolled {
		delta := cursor.Sub(c.scrollCursorStart)
		if stdmath.Abs(delta.X) > 0 || stdmath.Abs(delta.Y) > 0 {
			camera := archetype.MustFindCamera(w)
			cameraPos := transform.GetTransform(camera).LocalPosition
			cameraScale := transform.GetTransform(camera).LocalScale

			// Speed up when the camera is zoomed out
			delta.X /= cameraScale.X
			delta.Y /= cameraScale.Y

			newPos := cameraPos.Sub(delta)

			transform.GetTransform(camera).LocalPosition = engine.ClampVector(newPos, c.cameraBoundsMin, c.cameraBoundsMax)
			c.scrollCursorStart = cursor

			events.CameraMovedEvent.Publish(w, events.CameraMoved{
				Delta: delta,
			})
		}
	}
}

func (c *Controls) showCrosshair(
	w donburi.World,
	parent *donburi.Entry,
	pos component.Position,
) {
	crosshair := archetype.New(w).
		WithLayer(component.SpriteLayerCrosshair).
		WithSprite(component.SpriteData{

			Image: assets.Crosshair,
		}).
		With(component.Crosshair).
		Entry()

	component.Crosshair.SetValue(crosshair, component.CrosshairData{
		Position: pos,
	})

	transform.AppendChild(parent, crosshair, false)
}

func (c *Controls) selectTile(w donburi.World, entry *donburi.Entry) {
	c.deleteSelectable(w)

	selectable := component.Selectable.Get(entry)
	selectable.Selected = true

	archetype.New(w).
		WithParent(entry).
		WithLayer(component.SpriteLayerCrosshair).
		WithSprite(component.SpriteData{
			Image: assets.Selected,
		}).
		With(component.SelectedIndicator)

	events.TileSelectedEvent.Publish(w, events.TileSelected{
		Tile: entry,
	})
}

func (c *Controls) deleteSelectable(w donburi.World) {
	selected, ok := c.selectedQuery.First(w)
	if !ok {
		return
	}

	parent, ok := transform.GetParent(selected)
	if !ok {
		panic("selected indicator has no parent")
	}

	selectable := component.Selectable.Get(parent)
	selectable.Selected = false

	component.Destroy(selected)
}

func (c *Controls) unselect(w donburi.World) {
	c.deleteSelectable(w)

	selected, ok := c.selectedQuery.First(w)
	if !ok {
		return
	}

	parent, ok := transform.GetParent(selected)
	if !ok {
		panic("selected indicator has no parent")
	}

	events.TileUnselectedEvent.Publish(w, events.TileUnselected{
		Tile: parent,
	})
}

func distanceBetweenTouchIDs(touchID1, touchID2 ebiten.TouchID) float64 {
	x1, y1 := ebiten.TouchPosition(touchID1)
	x2, y2 := ebiten.TouchPosition(touchID2)
	return stdmath.Hypot(float64(x2-x1), float64(y2-y1))
}

func (c *Controls) Draw(w donburi.World, screen *ebiten.Image) {
	if !c.debug.Enabled {
		return
	}

	var scrollX, scrollY int
	if c.isTouchScrolling {
		scrollX, scrollY = ebiten.TouchPosition(c.scrollTouchID)
	}

	baseX := 32

	ebitenutil.DebugPrintAt(screen, "Controls", baseX, 480)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Touching: %v UI: %v (Timer: %v)", c.isTouching, c.isTouchingUI, c.touchTimer.PercentDone()), baseX, 500)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Touch scrolling: %v (%v, %v) touch ID: %v", c.isTouchScrolling, scrollX, scrollY, c.scrollTouchID), baseX, 520)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Map scrolled: %v", c.mapScrolled), baseX, 540)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pinching: %v (%v)", c.isPinching, c.zoomFactor), baseX, 560)
}

func spawnClickDebug(w donburi.World, cursorX int, cursorY int, clr color.Color) {
	imageSize := 32

	debugImage := ebiten.NewImage(32, 32)
	vector.StrokeLine(debugImage, 0, 0, 32, 32, 2, clr, true)
	vector.StrokeLine(debugImage, 0, 32, 32, 0, 2, clr, true)

	debug := archetype.New(w).
		WithPosition(math.Vec2{
			X: float64(cursorX - imageSize/2),
			Y: float64(cursorY - imageSize/2),
		}).
		WithLayer(component.SpriteUILayerTop).
		WithSprite(component.SpriteData{
			Image: debugImage,
		}).
		With(component.TimeToLive).
		With(component.UI).
		Entry()

	component.TimeToLive.SetValue(debug, component.TimeToLiveData{
		Timer: engine.NewTimer(2 * time.Second),
	})
}

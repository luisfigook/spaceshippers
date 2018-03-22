package main

import (
	"strconv"

	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type ShipDesignMenu struct {
	burl.BaseState

	window *burl.Container

	roomColumn        *burl.Container
	roomLists         *burl.PagedContainer
	installedRoomList *burl.List
	allRoomList       *burl.List
	roomDetails       *burl.Container

	shipView           *burl.TileView
	selectionAnimation *burl.PulseAnimation
	stars              StarField

	shipColumn *burl.Container

	buttons         *burl.Container
	addRemoveButton *burl.Button
	saveButton      *burl.Button
	loadButton      *burl.Button
	returnButton    *burl.Button
	helpText        *burl.Textbox

	dialog Dialog

	ship        *Ship
	roomToAdd   *Room
	roomAddGood bool

	roomTemplateOrder []RoomType //ordering for room templates, so we can later sort/filter

	offX, offY int //shipview camera offsets
}

func NewShipDesignMenu() (sdm *ShipDesignMenu) {
	sdm = new(ShipDesignMenu)

	sdm.window = burl.NewContainer(78, 43, 1, 1, 0, true)
	sdm.window.SetTitle("USE YOUR IMAGINATION")

	sdm.roomColumn = burl.NewContainer(20, 43, 0, 0, 0, true)
	sdm.roomColumn.Add(burl.NewTextbox(20, 1, 0, 0, 0, false, true, "Modules"))
	sdm.roomLists = burl.NewPagedContainer(20, 20, 0, 2, 0, true)

	sdm.allRoomList = burl.NewList(18, 16, 0, 0, 0, true, "No modules exists in the whole universe somehow.")
	sdm.installedRoomList = burl.NewList(18, 16, 0, 0, 0, true, "Ain't got no modules!")

	all := sdm.roomLists.AddPage("All")
	all.Add(sdm.allRoomList)
	installed := sdm.roomLists.AddPage("Installed")
	installed.Add(sdm.installedRoomList)

	sdm.roomDetails = burl.NewContainer(20, 20, 0, 23, 0, true)

	sdm.roomColumn.Add(sdm.roomLists, sdm.roomDetails)

	sdm.shipView = burl.NewTileView(36, 36, 21, 0, 0, false)
	sdm.shipColumn = burl.NewContainer(20, 43, 58, 0, 0, true)
	sdm.buttons = burl.NewContainer(36, 6, 21, 37, 0, true)

	sdm.addRemoveButton = burl.NewButton(8, 1, 0, 0, 0, true, true, "[A]dd Module")

	sdm.helpText = burl.NewTextbox(36, 4, 0, 2, 0, true, true, "")
	sdm.buttons.Add(sdm.addRemoveButton, sdm.helpText)

	sdm.window.Add(sdm.roomColumn, sdm.shipView, sdm.shipColumn, sdm.buttons)

	sdm.stars = NewStarField(20, sdm.shipView)

	sdm.ship = NewShip("whatever", nil)
	sdm.CenterView()

	sdm.roomTemplateOrder = make([]RoomType, 0)
	for i, _ := range roomTemplates {
		sdm.roomTemplateOrder = append(sdm.roomTemplateOrder, i)
	}

	sdm.UpdateAllRoomList()
	sdm.UpdateRoomDetails()
	sdm.UpdateHelpText()
	sdm.UpdateSelectionAnimation()

	return
}

func (sdm *ShipDesignMenu) CenterView() {
	w, h := sdm.shipView.Dims()
	if sdm.roomToAdd == nil {
		sdm.offX = sdm.ship.x + sdm.ship.width/2 - w/2
		sdm.offY = sdm.ship.y + sdm.ship.height/2 - h/2
	} else {
		sdm.offX = sdm.roomToAdd.X + sdm.roomToAdd.Width/2 - w/2
		sdm.offY = sdm.roomToAdd.Y + sdm.roomToAdd.Height/2 - h/2
	}
}

func (sdm *ShipDesignMenu) UpdateHelpText() {
	if sdm.roomToAdd != nil {
		sdm.helpText.ChangeText("ADDING MODULE: " + sdm.roomToAdd.Name + "/n/n Press ARROW KEYS to move, ENTER to add module to ship, and ESCAPE to cancel.")
	} else if sdm.roomLists.CurrentIndex() == 0 {
		sdm.helpText.ChangeText("Welcome to the Ship Designer!/n/n Use PGUP/PGDOWN to select a module to add. Press TAB to see all modules currently installed.")
	} else {
		sdm.helpText.ChangeText("Welcome to the Ship Designer!/n/n Use PGUP/PGDOWN to select a module to remove. Press TAB to see all available modules.")
	}
}

func (sdm *ShipDesignMenu) HandleKeypress(key sdl.Keycode) {
	if sdm.dialog != nil {
		sdm.dialog.HandleKeypress(key)
		return
	}

	if sdm.roomToAdd == nil {
		switch key {
		case sdl.K_a:
			if sdm.roomLists.CurrentIndex() == 0 {
				sdm.roomToAdd = CreateRoomFromTemplate(sdm.roomTemplateOrder[sdm.allRoomList.GetSelection()])
				sdm.roomToAdd.X = sdm.ship.shipMap.Width/2 - sdm.roomToAdd.Width/2
				sdm.roomToAdd.Y = sdm.ship.shipMap.Height/2 - sdm.roomToAdd.Height/2
				sdm.CenterView()
				sdm.UpdateRoomState()
				sdm.UpdateHelpText()
				sdm.addRemoveButton.Press()
			}
		case sdl.K_r:
			if sdm.roomLists.CurrentIndex() == 1 {
				room := sdm.ship.Rooms[sdm.installedRoomList.GetSelection()]
				sdm.ship.RemoveRoom(room)
				sdm.UpdateInstalledRoomList()
				sdm.addRemoveButton.Press()
			}
		case sdl.K_TAB:
			sdm.roomLists.HandleKeypress(key)
			sdm.UpdateRoomDetails()
			sdm.UpdateHelpText()
			sdm.UpdateSelectionAnimation()
			if sdm.roomLists.CurrentIndex() == 0 {
				sdm.addRemoveButton.ChangeText("[A]dd Module")
			} else {
				sdm.addRemoveButton.ChangeText("[R]emove Module")
			}

		case sdl.K_PAGEUP:
			if sdm.roomLists.CurrentIndex() == 0 {
				sdm.allRoomList.Prev()
			} else {
				sdm.installedRoomList.Prev()
				sdm.UpdateSelectionAnimation()
			}
			sdm.UpdateRoomDetails()
		case sdl.K_PAGEDOWN:
			if sdm.roomLists.CurrentIndex() == 0 {
				sdm.allRoomList.Next()
			} else {
				sdm.installedRoomList.Next()
				sdm.UpdateSelectionAnimation()
			}
			sdm.UpdateRoomDetails()
		}
	} else {
		switch key {
		case sdl.K_RETURN:
			sdm.AddRoomToShip()
			sdm.UpdateHelpText()
		case sdl.K_UP:
			sdm.roomToAdd.Y -= 1
			sdm.UpdateRoomState()
		case sdl.K_DOWN:
			sdm.roomToAdd.Y += 1
			sdm.UpdateRoomState()
		case sdl.K_LEFT:
			sdm.roomToAdd.X -= 1
			sdm.UpdateRoomState()
		case sdl.K_RIGHT:
			sdm.roomToAdd.X += 1
			sdm.UpdateRoomState()
		case sdl.K_ESCAPE:
			sdm.roomToAdd = nil
			sdm.UpdateHelpText()
		}
	}
}

func (sdm *ShipDesignMenu) UpdateSelectionAnimation() {
	sdm.shipView.RemoveAnimation(sdm.selectionAnimation)
	if len(sdm.installedRoomList.Elements) != 0 {
		room := sdm.ship.Rooms[sdm.installedRoomList.GetSelection()]
		sdm.selectionAnimation = burl.NewPulseAnimation(0, 0, room.Width, room.Height, 50, 0, true)
		sdm.shipView.AddAnimation(sdm.selectionAnimation)
		if sdm.roomLists.CurrentIndex() == 1 {
			sdm.selectionAnimation.Activate()
		}
	}
}

func (sdm *ShipDesignMenu) Update() {
	sdm.Tick++

	if sdm.Tick%10 == 0 {
		sdm.stars.Shift()
	}
}

func (sdm *ShipDesignMenu) UpdateRoomState() {
	if sdm.roomToAdd != nil {
		sdm.roomAddGood = sdm.ship.CheckRoomValidAdd(sdm.roomToAdd, sdm.roomToAdd.X, sdm.roomToAdd.Y)
	}
}

func (sdm *ShipDesignMenu) AddRoomToShip() {
	sdm.UpdateRoomState()
	if sdm.roomAddGood {
		sdm.ship.AddRoom(sdm.roomToAdd, sdm.roomToAdd.X, sdm.roomToAdd.Y)
		sdm.roomToAdd = nil
	}
	sdm.UpdateInstalledRoomList()
	sdm.UpdateRoomDetails()
}

func (sdm *ShipDesignMenu) UpdateInstalledRoomList() {
	sdm.installedRoomList.ClearElements()
	for _, r := range sdm.ship.Rooms {
		sdm.installedRoomList.Append(r.Name)
	}
	sdm.installedRoomList.CheckSelection()
	sdm.UpdateSelectionAnimation()
	sdm.UpdateRoomDetails()
}

func (sdm *ShipDesignMenu) UpdateAllRoomList() {
	sdm.allRoomList.ClearElements()
	for _, temp := range sdm.roomTemplateOrder {
		sdm.allRoomList.Append(roomTemplates[temp].name)
	}
}

func (sdm *ShipDesignMenu) UpdateRoomDetails() {
	sdm.roomDetails.ClearElements()

	var room *Room

	switch sdm.roomLists.CurrentIndex() {
	case 0: //All
		room = CreateRoomFromTemplate(sdm.roomTemplateOrder[sdm.allRoomList.GetSelection()])
	case 1: //Installed modules
		if len(sdm.installedRoomList.Elements) != 0 {
			room = sdm.ship.Rooms[sdm.installedRoomList.GetSelection()]
		}
	}

	if room != nil {
		sdm.roomDetails.Add(burl.NewTextbox(20, 1, 0, 0, 0, true, true, room.Name))
		sdm.roomDetails.Add(burl.NewTextbox(20, 2, 0, 2, 0, false, true, room.Description))
		sdm.roomDetails.Add(burl.NewTextbox(20, 1, 0, 5, 0, false, false, "Dims: ("+strconv.Itoa(room.Width)+"x"+strconv.Itoa(room.Height)+")"))
		sdm.roomDetails.Add(burl.NewTextbox(20, 1, 0, 7, 0, false, false, "STATS:"))
		for i, s := range room.Stats {
			sdm.roomDetails.Add(burl.NewTextbox(20, 1, 2, 8+i, 0, false, false, s.GetName()+": "+strconv.Itoa(s.Modifier)))
		}
	}
}

func (sdm *ShipDesignMenu) Render() {
	sdm.stars.Draw()
	sdm.ship.DrawToTileView(sdm.shipView, sdm.offX, sdm.offY)

	if sdm.roomToAdd != nil {
		for i := 0; i < sdm.roomToAdd.Width*sdm.roomToAdd.Height; i++ {

			x := i%sdm.roomToAdd.Width + sdm.roomToAdd.X - sdm.offX
			y := i/sdm.roomToAdd.Width + sdm.roomToAdd.Y - sdm.offY
			w, h := sdm.shipView.Dims()

			if burl.CheckBounds(x, y, w, h) {
				tv := sdm.roomToAdd.RoomMap.GetTile(i%sdm.roomToAdd.Width, i/sdm.roomToAdd.Width).GetVisuals()
				if sdm.roomAddGood {
					sdm.shipView.Draw(x, y, tv.Glyph, tv.ForeColour, burl.COL_GREEN)
				} else {
					sdm.shipView.Draw(x, y, tv.Glyph, tv.ForeColour, burl.COL_RED)
				}
			}
		}
	} else if sdm.roomLists.CurrentIndex() == 1 && len(sdm.installedRoomList.Elements) != 0 {
		room := sdm.ship.Rooms[sdm.installedRoomList.GetSelection()]
		sdm.selectionAnimation.MoveTo(room.X-sdm.offX, room.Y-sdm.offY)
	}

	sdm.window.Render()

	if sdm.dialog != nil {
		sdm.dialog.Render()
	}
}
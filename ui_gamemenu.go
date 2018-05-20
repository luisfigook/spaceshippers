package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type GameMenu struct {
	burl.PagedContainer

	playerPage *burl.Container

	//MISSIONS PAGE----------
	missionsPage           *burl.Container
	missionList            *burl.List
	missionDescriptionText *burl.Textbox
	missionStatusText      *burl.Textbox
	missionGoalList        *burl.List
	missionCriteriaList    *burl.List
	//-----------------------

	player *Player
}

func NewGameMenu(p *Player) (gm *GameMenu) {
	gm = new(GameMenu)
	gm.PagedContainer = *burl.NewPagedContainer(40, 28, 39, 3, 10, true)
	gm.SetVisibility(false)
	gm.SetHint("TAB to switch submenus")
	gm.player = p

	gm.playerPage = gm.AddPage("Player")

	gm.missionsPage = gm.AddPage("Missions")
	gm.missionList = burl.NewList(10, 23, 1, 1, 1, true, "No Missions To Do!")
	gm.missionList.SetHint("PgUp/PgDown")
	gm.missionStatusText = burl.NewTextbox(27, 1, 13, 0, 0, true, true, "")
	gm.missionDescriptionText = burl.NewTextbox(27, 4, 13, 2, 0, true, true, "")

	gm.missionGoalList = burl.NewList(27, 9, 13, 7, 0, true, "Nothing to do???")
	gm.missionGoalList.SetTitle("TO DO")
	gm.missionGoalList.ToggleHighlight()

	gm.missionCriteriaList = burl.NewList(27, 8, 13, 17, 0, true, "Nothing to do???")
	gm.missionCriteriaList.SetTitle("CRITERIA")
	gm.missionCriteriaList.ToggleHighlight()

	gm.missionsPage.Add(gm.missionList, gm.missionDescriptionText, gm.missionStatusText, gm.missionGoalList, gm.missionCriteriaList)

	return
}

func (gm *GameMenu) HandleKeypress(key sdl.Keycode) {
	gm.PagedContainer.HandleKeypress(key)

	switch gm.CurrentPage() {
	case gm.missionsPage:
		gm.missionList.HandleKeypress(key)
	}
}

func (gm *GameMenu) UpdateMissions() {
	gm.missionList.ClearElements()
	for _, mis := range gm.player.MissionLog {
		gm.missionList.Add(burl.NewTextbox(10, 2, 0, 0, 0, false, false, mis.name))
	}
	gm.missionList.Calibrate()

	if len(gm.player.MissionLog) != 0 {
		m := gm.player.MissionLog[gm.missionList.GetSelection()]

		switch m.status {
		case goal_COMPLETE:
			gm.missionStatusText.ChangeText("Mission Successful!")
		case goal_FAILED:
			gm.missionStatusText.ChangeText("Mission Failed!")
		case goal_INPROGRESS:
			gm.missionStatusText.ChangeText("Mission in Progress.")
		}

		gm.missionDescriptionText.ChangeText(m.description)

		gm.missionGoalList.ClearElements()
		for _, g := range m.steps {
			step := burl.NewContainer(27, 3, 0, 0, 0, false)
			step.Add(burl.NewTextbox(27, 1, 0, 0, 0, false, false, g.GetName()))
			step.Add(burl.NewTextbox(1, 1, 1, 1, 0, false, false, "- "))
			step.Add(burl.NewTextbox(25, 2, 2, 1, 0, false, false, g.GetDescription()))
			gm.missionGoalList.Add(step)
		}

		gm.missionCriteriaList.ClearElements()
		if len(m.criteria) == 0 {
			gm.missionCriteriaList.SetVisibility(false)
		} else {
			gm.missionCriteriaList.SetVisibility(true)
			for _, c := range m.criteria {
				criteria := burl.NewContainer(27, 3, 0, 0, 0, false)
				criteria.Add(burl.NewTextbox(27, 1, 0, 0, 0, false, false, c.GetName()))
				criteria.Add(burl.NewTextbox(26, 2, 1, 1, 0, false, false, "- "+c.GetDescription()))
				gm.missionCriteriaList.Add(criteria)
			}
		}
	} else {
		gm.missionDescriptionText.ChangeText("")
		gm.missionGoalList.ClearElements()
		gm.missionCriteriaList.ClearElements()
	}
}
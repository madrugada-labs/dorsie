package main

import (
	mario_go "github.com/madrugada-labs/mario-go"
	"github.com/rivo/tview"
)

type EventType int

const (
	GoToMainPage EventType = iota
	MarioBross
)

func StartEventLoop(app *tview.Application, screenManager *Screens, receiver <-chan EventType) {
	for event := range receiver {
		switch event {
		case GoToMainPage:
			mainPage := screenManager.GetMainPage()
			app.SetRoot(mainPage, true).SetFocus(mainPage)
			app.Sync()
			app.Draw()
		case MarioBross:
			app.Suspend(
				func() {
					mario_go.RunMarioGo()
				},
			)
		}
	}
}

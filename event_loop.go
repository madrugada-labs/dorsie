package main

import "github.com/rivo/tview"

type EventType int

const (
	GoToMainPage EventType = iota
)

func StartEventLoop(app *tview.Application, screenManager *Screens, receiver <-chan EventType) {
	for event := range receiver {
		switch event {
		case GoToMainPage:
			mainPage := screenManager.GetMainPage()
			app.SetRoot(mainPage, true).SetFocus(mainPage)
			app.Sync()
			app.Draw()
		}
	}
}

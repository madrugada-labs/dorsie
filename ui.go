package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func drawMainPage(app *tview.Application, dataFetcher *DataFetcher, filterSettings *FilterSettings, userPreferences *UserPreferences) {
	form := tview.NewForm().
		AddButton("See Jobs", func() {
			jobsView, err := dataFetcher.GetJobsPublic(filterSettings)
			if err != nil {
				log.Fatal(err)
			}
			app.SetRoot(jobsView, true).SetFocus(jobsView)
		}).
		AddButton("Config", func() {
			drawConfigPage(app, filterSettings, userPreferences)
		})
	app.SetRoot(form, true).SetFocus(form)
}

func drawConfigPage(app *tview.Application, filterSettings *FilterSettings, userPreferences *UserPreferences) {
	minSalary := filterSettings.MinSalary
	form := tview.NewForm().
		AddInputField("Min salary", fmt.Sprintf("%d", filterSettings.MinSalary), 20, nil, func(newMinSalary string) {
			s, err := strconv.Atoi(newMinSalary)
			if err != nil {
				return
			}
			minSalary = s
		}).
		AddButton("Save", func() {
			newState := UserPreferencesState{
				MinSalary: minSalary,
			}
			userPreferences.PersistPreferences(&newState)
			filterSettings.UpdateFilters(userPreferences)

			app.QueueEvent(tcell.NewEventKey(0, 'b', 0))
		}).
		AddButton("Back", func() {
			app.QueueEvent(tcell.NewEventKey(0, 'b', 0))
		})

	app.SetRoot(form, true).SetFocus(form)
}

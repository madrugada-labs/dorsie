package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func drawMainPage(app *tview.Application, dataFetcher *DataFetcher, filterSettings *FilterSettings, userPreferences *UserPreferences) {

	logo := tview.NewTextView().
		SetText(`
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@(    (@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&(((((@@@@@@
@@@@@@@@@@@@@@@@@@@(    (@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&(((((((((@@@@@@
@@@@@@@@@@@@@@@@@@@(    (@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&(((((((((((@@@@@@
@@@@@@@@@*        ((    (@@@@@@,          &@@@@@@     @,    @@@@           &@@@@@@@@,         .@@@@@@##(##&@@@@%##@@@@@@
@@@@@@@                 (@@@&                @@@@           @@     /@@@(     %@@@&      #@%      %@@%###@@@@@@@@@@%%%@@@
@@@@@@     %@@@@@@@     (@@%     @@@@@@@@     @@@      @@@@@@@     &@@@@@@@@@@@@&     @@@@@@@&    (@%%%%@@@@@@@@@###%@@@
@@@@@&     @@@@@@@@%    (@@     %@@@@@@@@.    &@@     @@@@@@@@&             (@@@                   @@@@##%@@@@@#####@@@@
@@@@@@     @@@@@@@@     (@@*    ,@@@@@@@@     @@@     @@@@@@@@@@@@@@@@&(      @@/    ,@@@@@@@@@@@@@@@@@###########&@@@@@
@@@@@@&      @@@@/      (@@@.     /@@@@      &@@@     @@@@@@@#     @@@@@@     @@@.     @@@@@(     @@@@@#########&@@@@@@@
@@@@@@@@*          &    (@@@@@.            %@@@@@     @@@@@@@@@.            /@@@@@@.            &@@@@@@#####&@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
	`).
		SetTextAlign(tview.AlignCenter)

	comment := tview.NewTextView().
		SetText(`browse jobs from the confort of your terminal ^.^`).
		SetTextAlign(tview.AlignCenter)
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
		}).SetButtonsAlign(1)

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(false).
		AddItem(logo, 1, 0, 5, 3, 0, 0, false).
		AddItem(comment, 6, 0, 4, 3, 0, 0, false).
		AddItem(form, 10, 0, 1, 3, 0, 0, true)

	app.SetRoot(grid, true).SetFocus(grid)
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

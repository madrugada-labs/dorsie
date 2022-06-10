package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
)

func drawMainPage(app *tview.Application, dataFetcher *DataFetcher, filterSettings *FilterSettings, userPreferences *UserPreferences) *tview.Grid {

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
			jobsPublic, err := dataFetcher.GetJobsPublic(filterSettings)
			jobsView := drawJobsView(app, jobsPublic)
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

	return grid
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

func drawJobsView(app *tview.Application, jobsPublic JobsPublic) *tview.Grid {
	// keep the original job list stored
	allJobsPublic := make(JobsPublic, len(jobsPublic))
	copy(allJobsPublic, jobsPublic)

	comment := tview.NewInputField().SetLabel("search: ").SetFieldWidth(40)
	jobsListUI := drawJobListUI(jobsPublic)

	grid := tview.NewGrid().
		SetRows(2).
		// SetColumns(30, 0, 30).
		SetBorders(false).
		AddItem(jobsListUI, 1, 0, 23, 3, 0, 0, true).
		AddItem(comment, 24, 0, 3, 3, 0, 0, false)

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// set 'q' to quit the app as well (aside from ctrl-c)
		switch event.Rune() {
		case rune(tcell.KeyTAB):
			event = nil
			app.SetFocus(comment)
		}
		return event
	})

	comment.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case rune(tcell.KeyEnter):
			grid.RemoveItem(jobsListUI)
			jobsPublic = filterJobsBasedOnSearch(comment.GetText(), allJobsPublic)
			jobsListUI = drawJobListUI(jobsPublic)
			grid.AddItem(jobsListUI, 1, 0, 23, 3, 0, 0, false)
			app.SetFocus(jobsListUI)
			event = nil
		}
		return event
	})

	return grid
}

func drawJobListUI(jobsPublic JobsPublic) *tview.List {
	jobsListUI := tview.NewList()
	for _, job := range jobsPublic {
		// copy the job.ID variable to avoid the GO gotcha :) in for loops: https://kkentzo.github.io/2021/01/21/golang-loop-variable-gotcha/
		jobID := job.ID.(string)
		jobsListUI = jobsListUI.AddItem(
			fmt.Sprintf("[::b] [%s - %s USD][-:-:-] %s @ %s", humanize.Comma(int64(job.MinSalary)), humanize.Comma(int64(job.MaxSalary)), job.Title, job.Company.Name),
			fmt.Sprintf("[green::]%s[-:-:-], [blue::]%s[-:-:-]", job.Experience, job.Field),
			'+',
			func() {
				openJob(jobID)
			},
		)
	}

	return jobsListUI
}

func filterJobsBasedOnSearch(search string, jobsPublic JobsPublic) JobsPublic {

	newJobsPublic := make(JobsPublic, 0)

	for _, job := range jobsPublic {
		if fuzzy.MatchNormalizedFold(search, string(job.Title)) || fuzzy.MatchNormalizedFold(search, string(job.Company.Name)) {
			newJobsPublic = append(newJobsPublic, job)
		}
	}

	return newJobsPublic

}

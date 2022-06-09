package main

import (
	"os"
	"os/signal"

	"github.com/gdamore/tcell/v2"
	"github.com/hasura/go-graphql-client"
	"github.com/rivo/tview"
)

// set the dorsie version, to ensure that the client is up to date
// with the server.
const dorsieVersion = "0.0.1"

func main() {
	// load some critical components
	userPreferences := NewUserPreferences()
	client := graphql.NewClient("https://persico.fly.dev/graphql", nil)
	dataFetcher := NewDataFetcher(client)

	err := userPreferences.CreatePreferencesFile()
	if err != nil {
		panic(err)
	}

	preferences, err := userPreferences.LoadPreferences()
	if err != nil {
		panic(err)
	}

	err = userPreferences.PersistPreferences(preferences)
	if err != nil {
		panic(err)
	}

	filterSettings := NewFilterSettings(userPreferences)

	// UI showtime
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// set 'q' to quit the app as well (aside from ctrl-c)
		switch event.Rune() {
		case 'q':
			event = nil
			app.Stop()
		case 'b':
			event = nil
			drawMainPage(app, dataFetcher, filterSettings, userPreferences)
		}
		return event
	})

	drawMainPage(app, dataFetcher, filterSettings, userPreferences)

	// crl-c handler - quit the app
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		app.Stop()
	}()

	// run!
	if err := app.Run(); err != nil {
		panic(err)
	}
}

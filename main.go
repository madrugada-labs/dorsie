package main

import (
	"os"
	"os/signal"

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
	mainPage := drawMainPage(app, dataFetcher, filterSettings, userPreferences)
	app.SetRoot(mainPage, true).SetFocus(mainPage)

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

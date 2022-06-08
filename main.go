package main

import (
	"log"
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
	client := graphql.NewClient("https://persico.fly.dev/graphql", nil)
	dataFetcher := NewDataFetcher(client)
	jobsView, err := dataFetcher.GetJobsPublic()
	if err != nil {
		log.Fatal(err)
	}

	app := tview.NewApplication()

	// UI showtime

	// crl-c handler - quit the app
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		app.Stop()
	}()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// set 'q' to quit the app as well (aside from ctrl-c)
		switch event.Rune() {
		case 'q':
			app.Stop()
		}
		return event
	})

	// run!
	if err := app.SetRoot(jobsView, true).SetFocus(jobsView).Run(); err != nil {
		panic(err)
	}
}

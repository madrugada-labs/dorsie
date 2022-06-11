package main

import (
	"os"
	"os/signal"
	"strings"

	"flag"

	"github.com/hasura/go-graphql-client"
	"github.com/rivo/tview"
)

// set the dorsie version, to ensure that the client is up to date
// with the server.
const dorsieVersion = "0.0.1"

var minSalary = flag.Int("minSalary", -1, "min salary for a role")
var fields = flag.String("fields", "", "fields of interest separated by comma: engineering,marketing")

type Flags struct {
	MinSalary *int
	Fields    []FieldEnum
}

func (f *Flags) UpdateFlags() {
	f.MinSalary = minSalary
	var fieldsArray []string
	if fields != nil && *fields != "" {
		fieldsArray = strings.Split(*fields, ",")
	}
	for _, field := range fieldsArray {
		f.Fields = append(f.Fields, FieldEnum(field))
	}
}

func main() {
	flag.Parse()

	flags := Flags{}
	flags.UpdateFlags()

	// load some critical components
	userPreferences := NewUserPreferences()
	client := graphql.NewClient("https://persico.fly.dev/graphql", nil)
	dataFetcher := NewDataFetcher(client)

	err := userPreferences.CreatePreferencesFile()
	if err != nil {
		panic(err)
	}

	preferences, err := userPreferences.LoadPreferences(flags)
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

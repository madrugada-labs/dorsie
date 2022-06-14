package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"flag"

	"github.com/hasura/go-graphql-client"
	"github.com/rivo/tview"
)

// set the dorsie version, to ensure that the client is up to date
// with the server.

var minSalary = flag.Int("minSalary", -1, "min salary for a role")
var fields = flag.String("fields", "", "fields of interest separated by comma: engineering,marketing")
var skipIntro = flag.Bool("skipIntro", false, "skip dorse's intro and go directly to jobs!")

type Flags struct {
	MinSalary *int
	Fields    []FieldEnum
	SkipIntro *bool
}

func (f *Flags) UpdateFlags() {
	f.MinSalary = minSalary
	f.SkipIntro = skipIntro
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

	// start event listener
	eventTypeChannel := make(chan EventType)

	screenManager := NewScreensManager()
	go StartEventLoop(app, screenManager, eventTypeChannel)

	screenManager.UpdateMainPage(drawMainPage(app, screenManager, eventTypeChannel, dataFetcher, filterSettings))

	jobsPublic, err := dataFetcher.GetJobsPublic(filterSettings)
	if err != nil {
		log.Fatal(err)
	}

	jobsView := drawJobsView(app, eventTypeChannel, jobsPublic)
	screenManager.UpdateJobView(jobsView)

	mainPage := screenManager.GetMainPage()

	if userPreferences.SkipIntroEnabled() {
		jobsViewUI := screenManager.GetJobsView()
		log.Println(jobsViewUI)
		app.SetRoot(jobsViewUI, true).SetFocus(jobsViewUI)
	} else {
		app.SetRoot(mainPage, true).SetFocus(mainPage)
	}

	// crl-c handler - quit the app
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.Stop()
	}()

	// run!
	if err := app.Run(); err != nil {
		panic(err)
	}
}

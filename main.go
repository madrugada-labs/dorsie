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
var skills = flag.String("skills", "", "skills related to the job by comma: rust,go,solidity,cloud")
var experiences = flag.String("experience", "early_career,mid_level,senior", "career experience separated by comma: early_career,mid_level,senior")
var skipIntro = flag.Bool("skipIntro", false, "skip dorse's intro and go directly to jobs!")

type Flags struct {
	MinSalary   *int
	Experiences []ExperienceEnum
	Fields      []FieldEnum
	Skills      []string
	SkipIntro   *bool
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
	var skillsArray []string
	if skills != nil && *skills != "" {
		skillsArray = strings.Split(*skills, ",")
	}
	f.Skills = append(f.Skills, skillsArray...)

	var experiencesArray []string
	if experiences != nil && *experiences != "" {
		experiencesArray = strings.Split(*experiences, ",")
	}
	for _, experience := range experiencesArray {
		f.Experiences = append(f.Experiences, ExperienceEnum(experience))
	}
	f.SkipIntro = skipIntro
	if !*f.SkipIntro {
		// if there's a flag enabled, we set skipIntro
		if *f.MinSalary != -1 || f.Fields != nil || f.Experiences != nil {
			t := true
			f.SkipIntro = &t
		}
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

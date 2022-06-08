package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/dustin/go-humanize"
	"github.com/hasura/go-graphql-client"
	"github.com/rivo/tview"
)

type DataFetcher struct {
	Client *graphql.Client
}

func NewDataFetcher(client *graphql.Client) *DataFetcher {
	return &DataFetcher{Client: client}
}

func (df *DataFetcher) GetJobsPublic() (*tview.List, error) {
	var query struct {
		JobsPublic JobsPublic
	}

	err := df.Client.Query(context.Background(), &query, nil)
	if err != nil {
		log.Println(err)
		// Handle error.
		return nil, err
	}
	jobsListUI := tview.NewList()
	for i, job := range query.JobsPublic {
		// copy the job.ID variable to avoid the GO gotcha :) in for loops: https://kkentzo.github.io/2021/01/21/golang-loop-variable-gotcha/
		jobID := job.ID.(string)
		jobsListUI = jobsListUI.AddItem(
			fmt.Sprintf("[::b] [%s - %s USD][-:-:-] %s @ %s", humanize.Comma(int64(job.MaxSalary)), humanize.Comma(int64(job.MaxSalary)), job.Title, job.Company.Name),
			fmt.Sprintf("[green::]%s[-:-:-], [blue::]%s[-:-:-]", job.Experience, job.Field),
			rune(i),
			func() {
				openJob(jobID)
			},
		)
	}
	return jobsListUI, nil

}

func openJob(jobAdID string) {
	url := fmt.Sprintf("https://dorse.io/job/%s", jobAdID)
	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("xdg-open", url).Start()
	case "darwin":
		_ = exec.Command("open", url).Start()
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		_ = fmt.Errorf("unsupported platform")
	}
}

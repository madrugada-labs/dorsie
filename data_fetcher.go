package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/hasura/go-graphql-client"
)

/// DataFetcher connects to Dorse's backend and retrieves data
type DataFetcher struct {
	Client *graphql.Client
}

/// NewDataFetcher is just a constructor for DataFetcher
func NewDataFetcher(client *graphql.Client) *DataFetcher {
	return &DataFetcher{Client: client}
}

/// GetJobsPublic hits the jobsPublic with several filters and returns the results (if available)
func (df *DataFetcher) GetJobsPublic(filterSettings *FilterSettings) (JobsPublic, error) {
	var query struct {
		JobsPublic JobsPublic `graphql:"jobsPublic(filters:{minSalary:$minSalary, fields:$fields, skills:$skills, experiences: $experiences})"`
	}

	skills := make([]graphql.String, 0)
	if filterSettings.Skills != nil {
		for _, skill := range *filterSettings.Skills {
			skills = append(skills, graphql.String(skill))
		}
	}
	variables := map[string]interface{}{
		"minSalary":   graphql.Int(filterSettings.MinSalary),
		"fields":      filterSettings.Fields,
		"skills":      skills,
		"experiences": filterSettings.Experiences,
	}

	err := df.Client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Println(err)
		// Handle error.
		return nil, err
	}
	return query.JobsPublic, nil

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

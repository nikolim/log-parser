package main

import (
	"fmt"
	"log"
	"log_analyser/api"
	"log_analyser/database"
	"log_analyser/types"
	"time"
)

func databaseError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func analyseRun(run types.Run, api api.GithubApi, db *database.Postgres) {
	db.StoreRun(run)
	jobs, steps, runners := api.GetJobs(run.Id)

	for _, runner := range runners {
		err := db.StoreRunner(runner)
		databaseError(err)
	}
	for _, job := range jobs {
		err := db.StoreJob(job)
		databaseError(err)
	}
	for _, step := range steps {
		err := db.StoreStep(step)
		databaseError(err)
	}
}

func analyseRepository(repository types.Repository, db *database.Postgres, config types.Config) {

	fmt.Println("Analysing repository: ", repository.Name.String, " from owner: ", repository.Owner.String)

	lastCrawl, err := db.GetLastCrawl(repository.Id)
	databaseError(err)
	githubApi := api.NewGithubAPI(repository.Owner.String, repository.Name.String, "token", config.GithubBaseUrl)

	runScanFrequency := time.Duration(config.RepoScanFrequency) * time.Minute
	currTime := time.Now()

	if currTime.Sub(lastCrawl) > runScanFrequency {
		err := db.UpdateLastCrawl(repository.Id, currTime)
		databaseError(err)
	}

	workflows := githubApi.GetWorkflows()
	for _, workflow := range workflows {
		workflow.RepositoryId = repository.Id
		err := db.StoreWorkflow(workflow)
		databaseError(err)
	}

	runs := githubApi.GetRuns("")
	for _, run := range runs {
		analyseRun(run, githubApi, db)
	}

}

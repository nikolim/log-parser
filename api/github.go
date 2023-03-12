package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log_analyser/types"
	"net/http"
	"net/url"
	"os"
	"time"
)

type GithubApi struct {
	Owner   string
	Repo    string
	Token   string
	BaseUrl string
	client  http.Client
}

func NewGithubAPI(owner, repo, token, baseUrl string) GithubApi {
	githubApi := GithubApi{
		Owner:   owner,
		Repo:    repo,
		Token:   os.Getenv("GITHUB_TOKEN"),
		BaseUrl: baseUrl,
	}
	githubApi.client = http.Client{
		Timeout: 10 * time.Second,
	}
	return githubApi
}

func (g GithubApi) getRequest(endpoint string) []byte {

	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.Token))

	resp, err := g.client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}

func (g GithubApi) GetWorkflows() []types.Workflow {

	endpoint := fmt.Sprintf("%s/repos/%s/%s/actions/workflows", g.BaseUrl, g.Owner, g.Repo)
	body := g.getRequest(endpoint)

	workflowResp := struct {
		Workflows []types.Workflow `json:"workflows"`
	}{}
	err := json.Unmarshal(body, &workflowResp)
	if err != nil {
		log.Fatalln(err)
	}
	return workflowResp.Workflows
}

// GetRuns dateString is the lower bound for date, format: YYYY-MM-DDTHH:MM:SSZ (e.g. 2022-12-24T12:30:00Z)
func (g GithubApi) GetRuns(dateString string) []types.Run {

	endpoint := fmt.Sprintf("%s/repos/%s/%s/actions/runs", g.BaseUrl, g.Owner, g.Repo)
	if len(dateString) > 0 {
		endpoint += "?created=%3E%3D" + url.QueryEscape(dateString)
	}
	body := g.getRequest(endpoint)
	runResp := struct {
		Runs []types.Run `json:"workflow_runs"`
	}{}
	err := json.Unmarshal(body, &runResp)
	if err != nil {
		log.Fatalln(err)
	}
	return runResp.Runs
}

func (g GithubApi) GetJobs(runId int) ([]types.Job, []types.Step, []types.Runner) {

	endpoint := fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/jobs", g.BaseUrl, g.Owner, g.Repo, runId)
	body := g.getRequest(endpoint)

	jobResp := struct {
		Jobs []types.Job `json:"jobs"`
	}{}
	err := json.Unmarshal(body, &jobResp)
	if err != nil {
		log.Fatalln(err)
	}

	jobs := jobResp.Jobs
	var steps []types.Step
	var runners []types.Runner
	runnerNames := make(map[int]bool)

	for _, job := range jobs {
		if _, ok := runnerNames[job.Runner.Id]; !ok {
			runnerNames[job.Runner.Id] = true
			runners = append(runners, job.Runner)
		}
		steps = append(steps, job.Steps...)
	}
	return jobs, steps, runners
}

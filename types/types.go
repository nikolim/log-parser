package types

import (
	"database/sql"
	"time"
)

type Runner struct {
	Id        int    `json:"runner_id"`
	Name      string `json:"runner_name"`
	GroupId   int    `json:"runner_group_id"`
	GroupName string `json:"runner_group_name"`
}

type Workflow struct {
	Id           int `json:"id"`
	RepositoryId string
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Run struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	RunNumber  string     `json:"run_id"`
	HeadBranch string     `json:"head_branch"`
	HeadSha    string     `json:"head_sha"`
	WorkflowId int        `json:"workflow_id"`
	Repository Repository `json:"repository"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt time.Time  `json:"finished_at"`
	Runner     Runner
	Steps      []Step `json:"steps"`
}

type Step struct {
	Id         string
	Name       string
	StartTime  time.Time
	EndTime    time.Time
	Conclusion string
	JobId      int
}

type Job struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"started_at"`
	EndTime    time.Time `json:"completed_at"`
	Conclusion string    `json:"conclusion"`
	JobId      int
	Steps      []Step `json:"steps"`
	Runner
}

type Repository struct {
	Id         string         `db:"id" toml:"id"`
	Name       sql.NullString `db:"name"`
	Owner      sql.NullString `db:"owner"`
	Platform   sql.NullString `db:"platform"`
	Collection sql.NullString `db:"collection"`
	Project    sql.NullString `db:"project"`
	LastCrawl  time.Time      `db:"last_crawl"`
}

type Config struct {
	RunFrequency      int    `toml:"RunFrequency"`
	RepoScanFrequency int    `toml:"RepoScanFrequency"`
	GithubBaseUrl     string `toml:"GithubBaseUrl"`
	PostgresHost      string `toml:"PostgresHost"`
	PostgresPort      int    `toml:"PostgresPort"`
	PostgresUser      string `toml:"PostgresUser"`
	PostgresPassword  string `toml:"PostgresPassword"`
	PostgresDb        string `toml:"PostgresDb"`
}

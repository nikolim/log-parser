package database

import (
	"database/sql"
	"fmt"
	"log"
	"log_analyser/types"
	"time"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(host, user, password, database string, port int) *Postgres {

	postgres := Postgres{}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, database, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// validate DB connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	postgres.db = db
	log.Println("Connected to Database")

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &postgres
}

func (p *Postgres) StoreRun(run types.Run) {
	query := `INSERT INTO runs (id, name, run_number, head_branch, head_sha, workflow_id, repository_id, started_at, finished_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
			  ON CONFLICT (id) DO UPDATE SET id=excluded.id 
			  returning id`
	_, err := p.db.Exec(query, run.Id, run.Name, run.RunNumber, run.HeadBranch, run.HeadSha, run.WorkflowId, run.Repository.Id, run.StartedAt, run.FinishedAt)
	if err != nil {
		log.Println(err)
	}
}

func (p *Postgres) GetRepositories() ([]types.Repository, error) {
	query := "SELECT * FROM repositories"
	rows, err := p.db.Query(query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var repos []types.Repository
	for rows.Next() {
		var repo types.Repository
		if err := rows.Scan(&repo.Id, &repo.Name, &repo.Owner, &repo.Platform, &repo.Collection, &repo.Project, &repo.LastCrawl); err != nil {
			return repos, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func (p *Postgres) StoreJob(job types.Job) error {
	query := `INSERT INTO jobs (id, name, start_time, end_time, runner, conclusion, run_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  ON CONFLICT (id) DO UPDATE SET id=excluded.id
			  returning id`
	_, err := p.db.Exec(query, job.Id, job.Name, job.StartTime, job.EndTime, job.Runner, job.Conclusion)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) StoreStep(step types.Step) error {
	query := `INSERT INTO steps (name, start_time, end_time, job_id, conclusion) 
			  VALUES ($1, $2, $3, $4, $5)
			  ON CONFLICT (id) DO UPDATE SET id=excluded.id
			  returning id`
	_, err := p.db.Exec(query, step.Name, step.StartTime, step.EndTime, step.JobId, step.Conclusion)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) StoreRunner(runner types.Runner) error {
	query := `INSERT INTO runners (id, name, group_id, group_name) 
			  VALUES ($1, $2, $3, $4)
			  ON CONFLICT (id) DO UPDATE SET id=excluded.id
			  returning id`
	_, err := p.db.Exec(query, runner.Id, runner.Name, runner.GroupId, runner.GroupName)
	return err
}

func (p *Postgres) GetLastCrawl(id string) (time.Time, error) {
	row := p.db.QueryRow("SELECT last_crawl FROM repositories WHERE id = $1", id)
	var lastCrawl time.Time
	err := row.Scan(&lastCrawl)
	return lastCrawl, err
}

func (p *Postgres) UpdateLastCrawl(id string, lastCrawl time.Time) error {
	query := "UPDATE repositories SET last_crawl = $1 WHERE id = $2"
	_, err := p.db.Exec(query, lastCrawl, id)
	return err
}

func (p *Postgres) StoreWorkflow(workflow types.Workflow) error {
	query := `INSERT INTO workflows (id, name, path, state, created_at, updated_at, repository_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  ON CONFLICT (id) DO UPDATE SET id=excluded.id
			  returning id`
	_, err := p.db.Exec(query, workflow.Id, workflow.Name, workflow.Path, workflow.State, workflow.CreatedAt, workflow.UpdatedAt, workflow.RepositoryId)
	return err
}

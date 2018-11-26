package pg_bouncer_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/LasTshaMAN/Go-Execute/jobs"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	db, err := sqlx.Open("postgres", "postgres://test:test@127.0.0.1:6543/test_db?sslmode=disable")
	require.NoError(t, err)
	db.SetMaxIdleConns(8)
	db.SetMaxOpenConns(8)

	const jobsAmount = 100
	executor := jobs.NewExecutor(20, jobsAmount)
	jobber := newJobber(db, jobsAmount)
	for i := 0; i < jobsAmount; i++ {
		job := jobber.newJob()
		err := executor.TryToEnqueue(job)
		require.NoError(t, err)

		fmt.Println("queued a job")

		// Main thread does something else, before it starts scheduling Postgres requests again
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("waiting for jobs to finish ...")
	err, failed := jobber.getResult()
	require.NoError(t, err)
	require.Zero(t, failed)

	require.NoError(t, db.Close())
}

type jobber struct {
	db *sqlx.DB
	results chan error
	jobsAmount int
}

func newJobber(db *sqlx.DB, maxJobsAmount int) *jobber {
	return &jobber{
		db: db,
		results: make(chan error, maxJobsAmount),
	}
}

func (f *jobber) newJob() func() {
	f.jobsAmount++
	if randBool() {
		// Simple job
		return func() {
			_, err := f.db.Exec("SELECT 1;")
			f.results <- err
		}
	}
	// Long-running job
	return func() {
		seconds := 1 + rand.Intn(5)
		query := fmt.Sprintf("SELECT pg_sleep(%d);", seconds)
		_, err := f.db.Exec(query)
		f.results <- err
	}
}

// err - error of the job that finished last
// failed - amount of failed jobs
func (f *jobber) getResult() (err error, failed int) {
	for i := 0; i < f.jobsAmount; i++ {
		err = <-f.results
		if err != nil {
			failed++
		}
	}
	close(f.results)
	return
}

func randBool() bool {
	return rand.Intn(2) == 0
}

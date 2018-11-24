package pg_bouncer_test

import (
	"fmt"
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

	executor := jobs.NewExecutor(20, 1)

	const jobsAmount = 1000

	res := make(chan error, jobsAmount)
	for i := 1; i <= jobsAmount; i++ {
		executor.Enqueue(func() {
			_, err := db.Exec("SELECT pg_sleep(1);")
			res <- err
		})
		fmt.Println("queued job")

		// Main thread does something else, before it starts scheduling Postgres requests again
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("waiting for jobs to finish ...")
	failed := 0
	for i := 1; i <= jobsAmount; i++ {
		err := <-res
		if err != nil {
			failed++
		}
	}

	close(res)
	require.NoError(t, db.Close())

	fmt.Printf("==> %d jobs failed\n", failed)
	require.Zero(t, failed)
}

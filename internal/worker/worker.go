// internal/worker/worker.go
package worker

import (
    "context"
    "log"
    "time"

    "job-processor/internal/job"
)

type Worker struct {
    ID int
    Store *job.Store
}

func (w Worker) Start(ctx context.Context, jobs <-chan job.Job) {
    for {
        select {
        case <-ctx.Done():
            log.Printf("Worker %d stopped\n", w.ID)
            return

        case j := <-jobs:
            w.process(ctx, j)
        }
    }
}

func (w Worker) process(ctx context.Context, j job.Job) {
    exampleJob := j.(job.ExampleJob)

	w.Store.Set(exampleJob.ID, job.StatusProcessing)
    retries := 3

    for i := 0; i < retries; i++ {
        err := j.Execute(ctx)
        if err == nil {
            w.Store.Set(exampleJob.ID, job.StatusDone)
            return
        }

        log.Printf("Worker %d retry %d\n", w.ID, i+1)
        time.Sleep(time.Duration(1<<i) * time.Second) // exponential backoff
    }
    w.Store.Set(exampleJob.ID, job.StatusFailed)
    log.Printf("Worker %d job failed permanently\n", w.ID)
}


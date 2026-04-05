
package queue

import "job-processor/internal/job"

type Queue struct {
    jobs chan job.Job
}

func New(size int) *Queue {
    return &Queue{
        jobs: make(chan job.Job, size),
    }
}

func (q *Queue) Push(j job.Job) {
    q.jobs <- j
}

func (q *Queue) Pop() <-chan job.Job {
    return q.jobs
}

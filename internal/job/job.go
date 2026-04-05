package job

import (
    "context"
    "fmt"
    "time"
)
type Job interface {
    Execute(ctx context.Context) error
}
type ExampleJob struct {
    ID int
}

func (j ExampleJob) Execute(ctx context.Context) error {
    fmt.Println("Processing job", j.ID)
    time.Sleep(time.Second)
    return nil
}
type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusDone       Status = "done"
	StatusFailed     Status = "failed"
)
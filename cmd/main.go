package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	myhttp "job-processor/internal/http"
	"job-processor/internal/job"
	"job-processor/internal/queue"
	"job-processor/internal/worker"
)

func main() {
	// context для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// очередь и store
	q := queue.New(100)
	store := job.NewStore()

	// worker pool
	for i := 0; i < 5; i++ {
		w := worker.Worker{
			ID:    i,
			Store: store,
		}
		go w.Start(ctx, q.Pop())
	}

	// HTTP handler
	handler := &myhttp.Handler{
		Queue: q,
		Store: store,
	}

	http.HandleFunc("/jobs", handler.CreateJob)
	http.HandleFunc("/jobs/status", handler.GetJob)

	// HTTP сервер
	go func() {
		log.Println("HTTP server started on :8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("shutting down...")

	cancel()
	time.Sleep(2 * time.Second)
}

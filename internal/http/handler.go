package http


import (
    "encoding/json"
    "net/http"
    "fmt"
    "job-processor/internal/job"
    "job-processor/internal/queue"
)

type Handler struct {
    Queue *queue.Queue
    Store *job.Store
}

type CreateJobRequest struct {
    ID int `json:"id"`
}

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req CreateJobRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	h.Store.Set(req.ID, job.StatusPending)

	j := job.ExampleJob{ID: req.ID}
	h.Queue.Push(j)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	var id int
	fmt.Sscanf(idStr, "%d", &id)

	status, ok := h.Store.Get(id)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": string(status),
	})
}
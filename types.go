package puppetmaster

import (
	"time"
)

// JobPagination holds information about the paginated jobs list.
type JobPagination struct {
	Jobs  []Job           `json:"data"`
}

// JobRequest defines how to create a job.
type JobRequest struct {
	Code    string            `json:"code"`
	Status  string            `json:"status"`
	Vars    map[string]string `json:"vars"`
	Modules map[string]string `json:"modules"`
}

// JobResponse is an api wrapper around a single job.
type JobResponse struct {
	Errors map[string][]string `json:"errors"`
	Data   Job                 `json:"data"`
}

// Job represents a complete job including status, results and logs.
type Job struct {
	UUID       string                 `json:"uuid"`
	Status     string                 `json:"status"`
	Code       string                 `json:"code"`
	Vars       map[string]string      `json:"vars"`
	Modules    map[string]string      `json:"modules"`
	Error      string                 `json:"error"`
	Logs       []Log                  `json:"logs"`
	Results    map[string]interface{} `json:"results"`
	StartedAt  *time.Time             `json:"started_at"`
	FinishedAt *time.Time             `json:"finished_at"`
	Duration   int                    `json:"duration"`
}

// A Log represents a log line yielded by the executor
type Log struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

package puppet_master

import (
	"time"
)

// PaginationLinks holds information about how to retrieve more / other entities.
type PaginationLinks struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}

// PaginationMeta informs about the overall pagination state. How many items there are in total, which page
// we're currently at and which items we got.
type PaginationMeta struct {
	Path        string `json:"path"`
	FirstPage   uint   `json:"first_page"`
	CurrentPage uint   `json:"current_page"`
	LastPage    uint   `json:"last_page"`
	PerPage     uint   `json:"per_page"`
	From        uint   `json:"from"`
	To          uint   `json:"to"`
	Total       uint   `json:"total"`
}

// JobPagination holds information about the paginated jobs list.
type JobPagination struct {
	Data  []Job           `json:"data"`
	Links PaginationLinks `json:"links"`
	Meta  PaginationMeta  `json:"meta"`
}

// JobRequest defines how to create a job.
type JobRequest struct {
	Code    string            `json:"code"`
	Status  string            `json:"status"`
	Vars    map[string]string `json:"vars"`
	Modules map[string]string `json:"modules"`
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

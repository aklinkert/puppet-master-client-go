package puppetmaster

import "errors"

// possible status values
const (
	StatusCreated = "created"
	StatusQueued  = "queued"
	StatusDone    = "done"
)

const authHeader = "Authorization"

var (
	// ErrEmptyAPIToken is thrown when the apiToken given to NewClient() is empty.
	ErrEmptyAPIToken = errors.New("apiToken may not be empty")

	// ErrNotFound is thrown when given job UUID was not found by the puppet master
	ErrNotFound = errors.New("job was not found by given UUID")

	// ErrEmptyCode is thrown when you try to create a job with empty code
	ErrEmptyCode = errors.New("given JobRequest's code may not be empty")
)



package puppetmaster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// Client represents a client to interact with the puppet-master API.
type Client struct {
	apiToken    string
	baseURL     *url.URL
	debug       bool
	syncSleepMs uint
}

// NewClient returns a new Client instance.
func NewClient(baseURL, apiToken string) (*Client, error) {
	apiToken = strings.TrimSpace(apiToken)
	if apiToken == "" {
		return nil, ErrEmptyAPIToken
	}

	c := &Client{
		apiToken:    apiToken,
		syncSleepMs: 500,
	}

	var err error
	if c.baseURL, err = url.Parse(baseURL); err != nil {
		return nil, err
	}

	return c, nil
}

// EnableDebugLogs enables debug logging of requests and responses.
func (c *Client) EnableDebugLogs() {
	c.debug = true
}

func (c *Client) addAuthentication(req *http.Request) {
	req.Header.Add(authHeader, fmt.Sprintf("Bearer %s", c.apiToken))
}

func (c *Client) buildURL(subPath string, query map[string]string) string {
	base := *c.baseURL
	base.Path = path.Join(base.Path, subPath)

	q := base.Query()
	for k, v := range query {
		q.Add(k, v)
	}

	base.RawQuery = q.Encode()

	return base.String()
}

// do sends a request and does additional logging of request and response, if debug mode is enabled.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if c.debug {
		dumpRequest(req)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.debug {
		dumpResponse(res)
	}

	return res, nil
}

// GetJobs returns all jobs as a paginated list
func (c *Client) GetJobs(page, perPage uint) (*JobPagination, error) {
	return c.GetJobsByStatus("", page, perPage)
}

// GetJobsByStatus lists jobs with the given status as a paginated list
func (c *Client) GetJobsByStatus(status string, page, perPage uint) (*JobPagination, error) {
	jobURL := c.buildURL("/jobs", map[string]string{
		"status":   status,
		"page":     string(page),
		"per_page": string(perPage),
	})

	req, err := http.NewRequest(http.MethodGet, jobURL, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthentication(req)

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, unexpectedResponse(res)
	}

	jobs := &JobPagination{}
	if err = json.NewDecoder(res.Body).Decode(jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// CreateJob schedules a new job for execution
func (c *Client) CreateJob(jobRequest *JobRequest) (*Job, error) {
	jobURL := c.buildURL("/jobs", map[string]string{})

	if strings.TrimSpace(jobRequest.Code) == "" {
		return nil, ErrEmptyCode
	}

	if jobRequest.Modules == nil {
		jobRequest.Modules = map[string]string{}
	}
	if jobRequest.Vars == nil {
		jobRequest.Vars = map[string]string{}
	}

	body, err := json.Marshal(jobRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, jobURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	c.addAuthentication(req)

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	job := &JobResponse{}
	if err = json.NewDecoder(res.Body).Decode(job); err != nil {
		return nil, err
	}

	if res.StatusCode != 201 {
		if res.StatusCode == 422 {
			return nil, unprocessableEntity(res, job.Errors)
		}

		return nil, unexpectedResponse(res)
	}

	return &job.Data, nil
}

// GetJob fetches a single job
func (c *Client) GetJob(uuid string) (*Job, error) {
	jobURL := c.buildURL(fmt.Sprintf("/jobs/%v", uuid), map[string]string{})
	req, err := http.NewRequest(http.MethodGet, jobURL, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthentication(req)

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			return nil, ErrNotFound
		}

		return nil, unexpectedResponse(res)
	}

	job := &JobResponse{}
	if err = json.NewDecoder(res.Body).Decode(job); err != nil {
		return nil, err
	}

	return &job.Data, nil
}

// DeleteJob deletes a job
func (c *Client) DeleteJob(uuid string) error {
	jobURL := c.buildURL(fmt.Sprintf("/jobs/%v", uuid), map[string]string{})
	req, err := http.NewRequest(http.MethodDelete, jobURL, nil)
	if err != nil {
		return err
	}

	c.addAuthentication(req)

	res, err := c.do(req)
	if err != nil {
		return err
	}

	if err := res.Body.Close(); err != nil {
		return fmt.Errorf("failed to close response body: %v", err)
	}

	if res.StatusCode != 204 {
		if res.StatusCode == 404 {
			return ErrNotFound
		}

		return unexpectedResponse(res)
	}

	return nil
}

// SetSyncSleepMs sets the amount of ms to sleep between two checks for the job being done
func (c *Client) SetSyncSleepMs(syncSleepMs uint) {
	c.syncSleepMs = syncSleepMs
}

// ExecuteSync executes a job synchronously by checking the job status in a changeable interval, 100ms by default.
// The amount of time is changeable by calling Client.SetSyncSleepMs(). Since the job is done even when an error occurred
// we can assure to return the finished job at some point in time.
func (c *Client) ExecuteSync(jobRequest *JobRequest) (*Job, error) {
	job, err := c.CreateJob(jobRequest)

	if err != nil {
		return nil, err
	}

	for {
		job, err = c.GetJob(job.UUID)
		if err != nil {
			if err == io.EOF {
				continue
			}

			return nil, err
		}

		if job.Status == StatusDone {
			return job, nil
		}

		time.Sleep(time.Duration(c.syncSleepMs) * time.Millisecond)
	}
}

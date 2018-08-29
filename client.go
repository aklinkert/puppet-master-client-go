package puppetmaster

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const authHeader = "Authorization"

var (
	// ErrEmptyAPIToken is thrown when the apiToken given to NewClient() is empty.
	ErrEmptyAPIToken = errors.New("apiToken may not be empty")

	// ErrNotFound is thrown when given job UUID was not found by the puppet master
	ErrNotFound = errors.New("job was not found by given UUID")
)

// Client represents a client to interact with the puppet-master API.
type Client struct {
	apiToken string
	baseURL  *url.URL
	debug    bool
}

// NewClient returns a new Client instance.
func NewClient(baseURL, apiToken string) (*Client, error) {
	apiToken = strings.TrimSpace(apiToken)
	if apiToken == "" {
		return nil, ErrEmptyAPIToken
	}

	c := &Client{
		apiToken: apiToken,
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

	if err == nil && c.debug {
		dumpResponse(res)
	}

	return res, err
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

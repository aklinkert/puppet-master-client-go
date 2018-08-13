package puppet_master

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
	ErrEmptyApiKey = errors.New("apiKey may not be empty")
)

// Client represents a client to interact with the puppet-master API.
type Client struct {
	apiKey  string
	baseURL *url.URL
}

// NewClient returns a new Client instance.
func NewClient(baseURL, apiKey string) (*Client, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, ErrEmptyApiKey
	}

	c := &Client{
		apiKey: apiKey,
	}

	var err error
	if c.baseURL, err = url.Parse(baseURL); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) addAuthentication(req *http.Request) {
	req.Header.Add(authHeader, fmt.Sprintf("Bearer %s", c.apiKey))
}

func (c *Client) buildUrl(subPath string, query map[string]string) string {
	base := *c.baseURL
	base.Path = path.Join(base.Path, subPath)

	q := base.Query()
	for k, v := range query {
		q.Add(k, v)
	}

	base.RawQuery = q.Encode()

	return base.String()
}

// GetJobs lists jobs with pagination
func (c *Client) GetJobs(status string, page, perPage uint) (*JobPagination, error) {
	jobUrl := c.buildUrl("/jobs", map[string]string{
		"status":   status,
		"page":     string(page),
		"per_page": string(perPage),
	})

	req, err := http.NewRequest(http.MethodGet, jobUrl, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthentication(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	jobs := &JobPagination{}
	if err = json.NewDecoder(res.Body).Decode(jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// GetJob fetches a single job
func (c *Client) GetJob(uuid string) (*Job, error) {
	jobUrl := c.buildUrl(fmt.Sprintf("/jobs/%v", uuid), map[string]string{})
	req, err := http.NewRequest(http.MethodGet, jobUrl, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthentication(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	job := &Job{}
	if err = json.NewDecoder(res.Body).Decode(job); err != nil {
		return nil, err
	}

	return job, nil
}

// CreateJob schedules a new job for execution
func (c *Client) CreateJob(jobRequest *JobRequest) (*Job, error) {
	jobUrl := c.buildUrl("/jobs", map[string]string{})

	body, err := json.Marshal(jobRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, jobUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	c.addAuthentication(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	job := &Job{}
	if err = json.NewDecoder(res.Body).Decode(job); err != nil {
		return nil, err
	}

	return job, nil
}

// DeleteJob deletes a job completely
func (c *Client) DeleteJob(jobRequest *JobRequest) (*Job, error) {
	jobUrl := c.buildUrl("/jobs", map[string]string{})

	body, err := json.Marshal(jobRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, jobUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	c.addAuthentication(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	job := &Job{}
	if err = json.NewDecoder(res.Body).Decode(job); err != nil {
		return nil, err
	}

	return job, nil
}

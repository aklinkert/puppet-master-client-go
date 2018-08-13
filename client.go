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
	// Thrown when the apiToken given to NewClient() is empty.
	ErrEmptyApiToken = errors.New("apiToken may not be empty")

	// Thrown when the teamSlug given to NewClient() is empty.
	ErrEmptyTeamSlug = errors.New("teamSlug may not be empty")

	// Thrown when given job UUID was not found by the puppet master
	ErrNotFound = errors.New("job was not found by given UUID")
)

// Client represents a client to interact with the puppet-master API.
type Client struct {
	apiToken, teamSlug string
	baseURL            *url.URL
	debug              bool
}

// NewClient returns a new Client instance.
func NewClient(teamSlug, baseURL, apiToken string) (*Client, error) {
	apiToken = strings.TrimSpace(apiToken)
	if apiToken == "" {
		return nil, ErrEmptyApiToken
	}

	teamSlug = strings.TrimSpace(teamSlug)
	if teamSlug == "" {
		return nil, ErrEmptyTeamSlug
	}

	c := &Client{
		apiToken: apiToken,
		teamSlug: teamSlug,
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

func (c *Client) buildUrl(subPath string, query map[string]string) string {
	base := *c.baseURL
	base.Path = path.Join(base.Path, "teams", c.teamSlug, subPath)

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

// GetAllJobs returns all jobs as a paginated list
func (c *Client) GetAllJobs(page, perPage uint) (*JobPagination, error) {
	return c.GetJobsByStatus("", page, perPage)
}

// GetJobsByStatus lists jobs with the given status as a paginated list
func (c *Client) GetJobsByStatus(status string, page, perPage uint) (*JobPagination, error) {
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
	jobUrl := c.buildUrl("/jobs", map[string]string{})

	body, err := json.Marshal(jobRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, jobUrl, bytes.NewReader(body))
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
	jobUrl := c.buildUrl(fmt.Sprintf("/jobs/%v", uuid), map[string]string{})
	req, err := http.NewRequest(http.MethodGet, jobUrl, nil)
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

// GetJob fetches a single job
func (c *Client) DeleteJob(uuid string) error {
	jobUrl := c.buildUrl(fmt.Sprintf("/jobs/%v", uuid), map[string]string{})
	req, err := http.NewRequest(http.MethodDelete, jobUrl, nil)
	if err != nil {
		return err
	}

	c.addAuthentication(req)

	res, err := c.do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 204 {
		if res.StatusCode == 404 {
			return ErrNotFound
		}

		return unexpectedResponse(res)
	}

	return nil
}

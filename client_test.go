package puppet_master

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddHeader(t *testing.T) {
	c := newTestClient(t, dumbHandler(200, nil))
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	t.Logf("Headers: %v", req.Header)
	c.client.addAuthentication(req)
	t.Logf("Headers: %v", req.Header)

	if len(req.Header) != 1 {
		t.Fatalf("Expected request to have 1 header, got %d", len(req.Header))
	}

	h := fmt.Sprintf("Bearer %v", c.apiToken)
	if req.Header.Get(authHeader) != h {
		t.Fatalf("Unexpected auth header %v, expected %v", req.Header.Get(authHeader), h)
	}
}

func TestBuildUrl(t *testing.T) {
	c := newTestClient(t, dumbHandler(200, nil))

	p := "jobs"
	q := map[string]string{
		"status":   "test",
		"per_page": "100",
	}

	exp := fmt.Sprintf("%s/teams/%s/jobs?per_page=100&status=test", c.endpoint, testTeamSlug)
	res := c.client.buildUrl(p, q)

	t.Logf("Build url %v, expected %v", res, exp)
	if res != exp {
		t.Errorf("Expected to get url %v, got %v", exp, res)
	}
}

func TestClient_CreateJob(t *testing.T) {
	res := readTestData(t, "create-response.json")
	c := newTestClient(t, dumbHandler(201, bytes.NewReader(res)))

	jobReq := &JobRequest{}
	readJSONFileInto(t, "create-request.json", jobReq)

	job, err := c.client.CreateJob(jobReq)
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	t.Logf("Created job: %+v", job)

	if job.UUID == "" {
		t.Errorf("Retrieved job has no UUID.")
	}

	if job.Code != jobReq.Code {
		t.Error("Job code differs")
	}
}

func TestClient_GetAllJobs(t *testing.T) {
	resData := readTestData(t, "get-jobs-response.json")
	c := newTestClient(t, dumbHandler(200, bytes.NewReader(resData)))

	res, err := c.client.GetAllJobs(1, 15)
	if err != nil {
		t.Fatalf("failed to list job: %v", err)
	}

	t.Logf("List number of jobs: %v", len(res.Jobs))

	if len(res.Jobs) != 10 {
		t.Errorf("Expected GetAllJobs to return 10 jobs, got %d", len(res.Jobs))
	}

	if res.Meta.Total != 10 {
		t.Errorf("Expected GetAllJobs to a total of 10, got %d", res.Meta.Total)
	}
}

func TestClient_GetJob(t *testing.T) {
	resData := readTestData(t, "get-job-response.json")
	c := newTestClient(t, dumbHandler(200, bytes.NewReader(resData)))

	uuid := "73e3a9b5-81c8-4743-9a7e-e80474c1b6e3"
	job, err := c.client.GetJob(uuid)
	if err != nil {
		t.Fatalf("failed to get job: %v", err)
	}

	if job == nil {
		t.Fatal("Returned job is nil")
	}

	if job.UUID != uuid {
		t.Errorf("Expected to get UUID %v, got %v", uuid, job.UUID)
	}
}

func TestClient_GetJobNotFound(t *testing.T) {
	c := newTestClient(t, dumbHandler(404, nil))

	uuid := "73e3a9b5-81c8-4743-9a7e-e80474c1b6e3"
	_, err := c.client.GetJob(uuid)
	if err != ErrNotFound {
		t.Fatalf("Expected to get ErrNotFound, got %v", err)
	}
}

func TestClient_DeleteJob(t *testing.T) {
	c := newTestClient(t, dumbHandler(204, nil))

	uuid := "73e3a9b5-81c8-4743-9a7e-e80474c1b6e3"
	err := c.client.DeleteJob(uuid)
	if err != nil {
		t.Fatalf("failed to delete job: %v", err)
	}
}

func TestClient_DeleteJobNotFound(t *testing.T) {
	c := newTestClient(t, dumbHandler(404, nil))

	uuid := "73e3a9b5-81c8-4743-9a7e-e80474c1b6e3"
	err := c.client.DeleteJob(uuid)
	if err != ErrNotFound {
		t.Fatalf("Expected to get ErrNotFound, got %v", err)
	}
}

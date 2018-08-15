package puppetmaster

import (
	"testing"
	"time"
)

func TestEqual(t *testing.T) {
	now := time.Now()
	cases := []struct {
		j1, j2 *Job
		equal  bool
	}{
		{
			j1:    &Job{},
			j2:    &Job{},
			equal: true,
		},
		{
			j1:    &Job{Status: "test"},
			j2:    &Job{Status: "test"},
			equal: true,
		},
		{
			j1:    &Job{Status: "test1"},
			j2:    &Job{Status: "test2"},
			equal: false,
		},
		{
			j1:    &Job{Code: "test"},
			j2:    &Job{Code: "test"},
			equal: true,
		},
		{
			j1:    &Job{Code: "test1"},
			j2:    &Job{Code: "test2"},
			equal: false,
		},
		{
			j1:    &Job{Error: "test"},
			j2:    &Job{Error: "test"},
			equal: true,
		},
		{
			j1:    &Job{Error: "test1"},
			j2:    &Job{Error: "test2"},
			equal: false,
		},
		{
			j1:    &Job{UUID: "test"},
			j2:    &Job{UUID: "test"},
			equal: true,
		},
		{
			j1:    &Job{UUID: "test1"},
			j2:    &Job{UUID: "test2"},
			equal: false,
		},
		{
			j1:    &Job{Vars: map[string]string{"test": "yeah"}},
			j2:    &Job{Vars: map[string]string{"test": "yeah"}},
			equal: true,
		},
		{
			j1:    &Job{Vars: map[string]string{"test1": "yeah"}},
			j2:    &Job{Vars: map[string]string{"test2": "yeah"}},
			equal: false,
		},
		{
			j1:    &Job{Vars: map[string]string{"test": "yeah1"}},
			j2:    &Job{Vars: map[string]string{"test": "yeah2"}},
			equal: false,
		},
		{
			j1:    &Job{Modules: map[string]string{"test": "yeah"}},
			j2:    &Job{Modules: map[string]string{"test": "yeah"}},
			equal: true,
		},
		{
			j1:    &Job{Modules: map[string]string{"test1": "yeah"}},
			j2:    &Job{Modules: map[string]string{"test2": "yeah"}},
			equal: false,
		},
		{
			j1:    &Job{Modules: map[string]string{"test": "yeah1"}},
			j2:    &Job{Modules: map[string]string{"test": "yeah2"}},
			equal: false,
		},
		{
			j1:    &Job{Results: map[string]interface{}{"test": "yeah"}},
			j2:    &Job{Results: map[string]interface{}{"test": "yeah"}},
			equal: true,
		},
		{
			j1:    &Job{Results: map[string]interface{}{"test1": "yeah"}},
			j2:    &Job{Results: map[string]interface{}{"test2": "yeah"}},
			equal: false,
		},
		{
			j1:    &Job{Results: map[string]interface{}{"test": "yeah1"}},
			j2:    &Job{Results: map[string]interface{}{"test": "yeah2"}},
			equal: false,
		},
		{
			j1:    &Job{Logs: []Log{}},
			j2:    &Job{Logs: []Log{}},
			equal: true,
		},
		{
			j1:    &Job{Logs: []Log{{}}},
			j2:    &Job{Logs: []Log{{}}},
			equal: true,
		},
		{
			j1:    &Job{Logs: []Log{{Level: "test"}}},
			j2:    &Job{Logs: []Log{{Level: "test"}}},
			equal: true,
		},
		{
			j1:    &Job{Logs: []Log{{Message: "test1"}}},
			j2:    &Job{Logs: []Log{{Message: "test2"}}},
			equal: false,
		},
		{
			j1:    &Job{Logs: []Log{{Time: now}}},
			j2:    &Job{Logs: []Log{{Time: now}}},
			equal: true,
		},
		{
			j1:    &Job{Logs: []Log{{Time: now}}},
			j2:    &Job{Logs: []Log{{Time: now.Add(1 * time.Second)}}},
			equal: false,
		},
		{
			j1:    &Job{Logs: []Log{{Time: now}, {Time: now}}},
			j2:    &Job{Logs: []Log{{Time: now}, {Time: now}}},
			equal: true,
		},
		{
			j1:    &Job{Logs: []Log{{Time: now}, {Time: now}}},
			j2:    &Job{Logs: []Log{{Time: now}, {Time: now.Add(1 * time.Second)}}},
			equal: false,
		},
	}

	for i, c := range cases {
		res1 := c.j1.Equal(c.j2)
		res2 := c.j2.Equal(c.j1)

		t.Logf("Case %d, res1=%v res2=%v, exp=%v", i, res1, res2, c.equal)
		if res1 != c.equal {
			t.Errorf("Case %d: Expected j1.Equal(j2) == %v, but got %v", i, c.equal, res1)
		}

		if res2 != c.equal {
			t.Errorf("Case %d: Expected j2.Equal(j1) == %v, but got %v", i, c.equal, res2)
		}
	}
}

func TestDatesAreEqual(t *testing.T) {
	now1 := time.Now()
	now2 := time.Now().Add(-1 * time.Hour)

	cases := []struct {
		t1, t2 *time.Time
		equal  bool
	}{
		{t1: nil, t2: nil, equal: true},
		{t1: &now1, t2: &now1, equal: true},
		{t1: nil, t2: &now1, equal: false},
		{t1: &now1, t2: nil, equal: false},
		{t1: &now1, t2: &now2, equal: false},
		{t1: &now2, t2: &now1, equal: false},
		{t1: &now2, t2: &now2, equal: true},
	}

	for i, c := range cases {
		res := datesAreEqual(c.t1, c.t2)

		t.Logf("case %d, res=%v, exp=%v", i, res, c.equal)

		if res != c.equal {
			t.Errorf("case %d: Expected equality=%v, got %v", i, c.equal, res)
		}
	}
}

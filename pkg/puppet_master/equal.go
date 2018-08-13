package puppet_master

import (
	"reflect"
	"time"
)

// Equal returns true when both given Jobs are equal
func (j *Job) Equal(j2 *Job) bool {
	return j.UUID == j2.UUID &&
		j.Status == j2.Status &&
		j.Code == j2.Code &&
		reflect.DeepEqual(j.Modules, j2.Modules) &&
		reflect.DeepEqual(j.Vars, j2.Vars) &&
		reflect.DeepEqual(j.Results, j2.Results) &&
		reflect.DeepEqual(j.Logs, j2.Logs) &&
		datesAreEqual(j.StartedAt, j2.StartedAt) &&
		datesAreEqual(j.FinishedAt, j2.FinishedAt) &&
		j.Error == j2.Error
}

func datesAreEqual(t1 *time.Time, t2 *time.Time) bool {
	if (t1 == nil && t2 != nil) || (t1 != nil && t2 == nil) {
		return false
	}

	if t1 == nil && t2 == nil {
		return true
	}

	return (*t1).Equal(*t2)
}

package tracing

import (
	"gitlab.com/akita/akita"
)

// A TaskStep represents a milestone in the processing of task
type TaskStep struct {
	Time akita.VTimeInSec `json:"time"`
	What string           `json:"what"`
}

// A Task is a task
type Task struct {
	ID        string           `json:"id"`
	ParentID  string           `json:"parent_id"`
	Kind      string           `json:"kind"`
	What      string           `json:"what"`
	Where     string           `json:"where"`
	StartTime akita.VTimeInSec `json:"start_time"`
	EndTime   akita.VTimeInSec `json:"end_time"`
	Steps     []TaskStep       `json:"steps"`
	Detail    interface{}      `json:"-"`
}

// TaskFilter is a function that can filter interesting tasks. If this function
// returns true, the task is considered useful.
type TaskFilter func(t Task) bool

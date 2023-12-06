package job

import (
	"time"
)

type State struct {
	Status Status `json:"status,omitempty"`

	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`

	StartTime    time.Time  `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	ObservedTime time.Time  `json:"observed_time,omitempty"`

	ExitCode *int `json:"exit_code,omitempty"`
}

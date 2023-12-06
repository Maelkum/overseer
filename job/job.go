package job

import (
	"io"
)

type Job struct {
	ID           string    `json:"id,omitempty"`
	Exec         Command   `json:"exec,omitempty"`
	Stdin        io.Reader `json:"stdin,omitempty"`
	OutputStream string    `json:"output_stream,omitempty"`
	Files        string    `json:"files,omitempty"`
	ErrorStream  string    `json:"error_stream,omitempty"`
	Limits       *Limits   `json:"job_limits,omitempty"`
}

type Command struct {
	Path string   `json:"path,omitempty"`
	Args []string `json:"args,omitempty"`
	Env  []string `json:"env,omitempty"`
}

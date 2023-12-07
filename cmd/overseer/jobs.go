package main

import (
	"github.com/Maelkum/overseer/job"
)

func getJobs() []job.Job {

	var jobs []job.Job

	srv1 := job.Job{
		Exec: job.Command{
			Path: `/home/aco/code/Maelkum/overseer/cmd/example-server/example-server`,
			Args: []string{
				"--address",
				":8080",
				"--name",
				"first-server-name",
				"--exec",
				"true",
			},
		},
		OutputStream: "http://localhost:9000/",
		ErrorStream:  "http://localhost:9001/",
		Limits: &job.Limits{
			NoExec:        true,
			CPUPercentage: 0.80,
			MemoryLimitKB: 128_000,
		},
	}

	srv2 := job.Job{
		Exec: job.Command{
			Path: `/home/aco/code/Maelkum/overseer/cmd/example-server/example-server`,
			Args: []string{
				"--address",
				":8081",
				"--name",
				"second-server-name"},
		},
		Limits: &job.Limits{
			CPUPercentage: 0.75,
			MemoryLimitKB: 256_000,
		},
	}

	jobs = append(jobs, srv1, srv2)

	return jobs
}

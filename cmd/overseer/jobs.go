package main

import (
	"github.com/google/uuid"

	"github.com/Maelkum/overseer/overseer"
)

func getJobs() []overseer.Job {

	var jobs []overseer.Job

	srv1 := overseer.Job{
		Exec: overseer.Command{
			Path: `/home/aco/code/Maelkum/overseer/cmd/example-server/example-server`,
			Args: []string{
				"--address",
				":8080",
				"--name",
				"first-server-name",
				"--execute",
				"true",
			},
		},
		ID:           uuid.New().String(),
		OutputStream: "http://localhost:9000/",
		ErrorStream:  "http://localhost:9001/",
		Limits: &overseer.JobLimits{
			NoExec: true,
		},
	}

	srv2 := overseer.Job{
		Exec: overseer.Command{
			Path: `/home/aco/code/Maelkum/overseer/cmd/example-server/example-server`,
			Args: []string{
				"--address",
				":8081",
				"--name",
				"second-server-name"},
		},
		ID: uuid.New().String(),
	}

	jobs = append(jobs, srv1, srv2)

	return jobs
}

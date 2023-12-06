package overseer

import (
	"errors"
	"fmt"
	"time"

	"github.com/Maelkum/overseer/job"
)

func (o *Overseer) Wait(id string) (job.State, error) {

	o.Lock()
	h, ok := o.jobs[id]
	o.Unlock()

	if !ok {
		return job.State{}, errors.New("unknown job")
	}

	h.Lock()
	defer h.Unlock()

	defer o.harvest(id)

	err := h.cmd.Wait()
	if err != nil {
		return job.State{}, fmt.Errorf("could not wait on job: %w", err)
	}

	endTime := time.Now()

	state := job.State{
		Status:       job.StatusDone,
		Stdout:       h.stdout.String(),
		Stderr:       h.stderr.String(),
		StartTime:    h.start,
		EndTime:      &endTime,
		ObservedTime: time.Now(),
	}

	exitCode := h.cmd.ProcessState.ExitCode()
	state.ExitCode = &exitCode
	if *state.ExitCode != 0 {
		state.Status = job.StatusFailed
	}

	err = o.limiter.DeleteGroup(id)
	if err != nil {
		o.log.Error().Err(err).Str("job", id).Msg("could not delete limit group")
	}

	return state, nil
}

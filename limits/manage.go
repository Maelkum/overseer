package limits

import (
	"errors"
	"fmt"
)

func (l *Limiter) AssignToGroup(name string, pid uint64) error {

	l.Lock()
	defer l.Unlock()

	cg, ok := l.limits[name]
	if !ok {
		return errors.New("unknown group")
	}

	l.log.Info().Str("name", name).Uint64("pid", pid).Msg("assigning process to limit group")

	err := cg.AddProc(pid)
	if err != nil {
		return fmt.Errorf("could not assign process to group: %w", err)
	}

	l.log.Info().Str("name", name).Uint64("pid", pid).Msg("process assigned to limit group")

	return nil
}

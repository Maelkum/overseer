package limits

import (
	"fmt"
	"path"

	"github.com/containerd/cgroups/v3/cgroup2"
)

func (l *Limiter) Create(name string, limits Limits) error {

	l.Lock()
	defer l.Unlock()

	_, ok := l.limits[name]
	if ok {
		return fmt.Errorf("limits with id %v already exist", name)
	}

	group := path.Join(l.mountpoint, l.cgroup, name)
	specs := limitsToResources(limits)

	cg, err := cgroup2.NewManager(l.mountpoint, group, specs)
	if err != nil {
		return fmt.Errorf("could not create cgroup: %w", err)
	}

	l.limits[name] = cg

	return nil
}

func (l *Limiter) Delete(name string) error {

	l.Lock()
	defer l.Unlock()

	cg, ok := l.limits[name]
	if !ok {
		return nil
	}

	err := cg.Delete()
	if err != nil {
		return fmt.Errorf("could not delete cgroup: %w", err)
	}

	delete(l.limits, name)

	return nil
}

package limits

import (
	"fmt"
	"os"
	"path"

	"github.com/containerd/cgroups/v3/cgroup2"
)

func (l *Limiter) CreateGroup(name string, limits Limits) error {

	l.Lock()
	defer l.Unlock()

	_, ok := l.limits[name]
	if ok {
		return fmt.Errorf("limits with id %v already exist", name)
	}

	l.log.Info().Str("name", name).Msg("creating limit group")

	group := path.Join(l.cgroup, name)
	specs := limitsToResources(limits)

	cg, err := cgroup2.NewManager(l.mountpoint, group, specs)
	if err != nil {
		return fmt.Errorf("could not create cgroup (cgroup: %v): %w", group, err)
	}

	l.limits[name] = cg

	l.log.Info().Str("name", name).Msg("limit group created")

	return nil
}

func (l *Limiter) ListGroups() ([]string, error) {

	path := path.Join(l.mountpoint, l.cgroup)

	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open limiter root cgroup (path: %v): %w", path, err)
	}

	entries, err := dir.ReadDir(0)
	if err != nil {
		return nil, fmt.Errorf("could not list limiter root cgroup: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		names = append(names, entry.Name())
	}

	return names, nil
}

func (l *Limiter) DeleteGroup(name string) error {

	l.Lock()
	defer l.Unlock()

	cg, ok := l.limits[name]
	if !ok {
		return nil
	}

	l.log.Info().Str("name", name).Msg("deleting limit group")

	err := cg.Delete()
	if err != nil {
		return fmt.Errorf("could not delete cgroup: %w", err)
	}

	delete(l.limits, name)

	l.log.Info().Str("name", name).Msg("limit group deleted")

	return nil
}

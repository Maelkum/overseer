package limits

import (
	"errors"
	"sync"

	"github.com/containerd/cgroups/v3"
	"github.com/containerd/cgroups/v3/cgroup2"
)

type Limiter struct {
	*sync.Mutex

	mountpoint string
	cgroup     string

	limits map[string]*cgroup2.Manager
}

func New(mountpoint string, parentCgroup string) (*Limiter, error) {

	// Check if the system supports cgroups v2.
	var haveV2 bool
	if cgroups.Mode() == cgroups.Unified {
		haveV2 = true
	}
	if !haveV2 {
		return nil, errors.New("cgroups v2 is not supported")
	}

	l := Limiter{
		mountpoint: mountpoint,
		cgroup:     parentCgroup,

		Mutex:  &sync.Mutex{},
		limits: make(map[string]*cgroup2.Manager),
	}

	return &l, nil
}

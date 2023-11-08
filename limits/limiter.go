package limits

import (
	"errors"
	"os"
	"sync"

	"github.com/containerd/cgroups/v3"
	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/rs/zerolog"
)

type Limiter struct {
	*sync.Mutex
	log zerolog.Logger

	mountpoint string
	cgroup     string

	limits map[string]*limitHandler
}

type limitHandler struct {
	manager *cgroup2.Manager
	handle  *os.File
}

func New(log zerolog.Logger, mountpoint string, parentCgroup string) (*Limiter, error) {

	// Check if the system supports cgroups v2.
	var haveV2 bool
	if cgroups.Mode() == cgroups.Unified {
		haveV2 = true
	}
	if !haveV2 {
		return nil, errors.New("cgroups v2 is not supported")
	}

	l := Limiter{
		log: log,

		mountpoint: mountpoint,
		cgroup:     parentCgroup,

		Mutex:  &sync.Mutex{},
		limits: make(map[string]*limitHandler),
	}

	l.log.Debug().Str("cgroup", l.cgroup).Msg("created limiter")

	return &l, nil
}

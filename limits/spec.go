package limits

import (
	"time"

	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func limitsToResources(limits Limits) *cgroup2.Resources {

	lr := specs.LinuxResources{}

	if limits.CPUPercentage != 1.0 {

		// We want to set total CPU time limit. We'll use one year as the period.
		period := uint64(time.Second.Microseconds())
		quota := int64(float64(period) * float64(limits.CPUPercentage))

		lr.CPU = &specs.LinuxCPU{
			Period: &period,
			Quota:  &quota,
		}
	}

	if limits.MemoryKB > 0 {

		// Convert limit to bytes.
		memLimit := limits.MemoryKB * 1000

		lr.Memory = &specs.LinuxMemory{
			Limit: &memLimit,
		}
	}

	return cgroup2.ToResources(&lr)
}

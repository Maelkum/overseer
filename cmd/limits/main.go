package main

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/Maelkum/overseer/limits"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {

	var (
		flagCgroup string

		flagName          string
		flagCPUPercentage float64
		flagMemoryKB      int64
	)

	// TODO: Consider sending a PR to remove the requirement to have a preceeding '/'.
	pflag.StringVar(&flagCgroup, "cgroup", "overseer", "default cgroup to use as parent for other settings")
	pflag.StringVar(&flagName, "name", "", "name to use for resource group")
	pflag.Float64Var(&flagCPUPercentage, "cpu-limit", 1.00, "CPU percentage to set")
	pflag.Int64Var(&flagMemoryKB, "memory-limit-kb", -1, "memory limit to set (in KB)")

	pflag.Parse()

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	if flagName == "" {
		log.Info().Msg("name not specified, setting the value for the root cgroup")
	}

	if !strings.HasPrefix(flagCgroup, "/") {
		flagCgroup = "/" + flagCgroup
	}

	limiter, err := limits.New(limits.DefaultMountpoint, flagCgroup)
	if err != nil {
		log.Error().Err(err).Msg("could not create limiter")
		return failure
	}

	limits := limits.Limits{
		CPUPercentage: flagCPUPercentage,
		MemoryKB:      flagMemoryKB,
	}
	err = limiter.Create(flagName, limits)
	if err != nil {
		log.Error().Err(err).Msg("could not set resource limits")
		return failure
	}

	log.Info().Msg("all done")

	return success
}

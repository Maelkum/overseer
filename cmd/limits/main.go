package main

import (
	"os"

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
		flagCgroup        string
		flagCPUPercentage float64
		flagMemoryKB      int64
	)

	pflag.StringVar(&flagCgroup, "cgroup", "overseer", "default cgroup to use as parent for other settings")
	pflag.Float64Var(&flagCPUPercentage, "cpu-limit", 1.00, "CPU percentage to set")
	pflag.Int64Var(&flagMemoryKB, "memory-limit-kb", -1, "memory limit to set (in KB)")

	pflag.Parse()

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	limiter, err := limits.New(limits.DefaultMountpoint, flagCgroup)
	if err != nil {
		log.Error().Err(err).Msg("could not create limiter")
		return failure
	}

	_ = limiter

	log.Info().Msg("all done")

	return success
}

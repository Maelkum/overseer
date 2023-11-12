package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"

	"github.com/Maelkum/overseer/limits"
	"github.com/Maelkum/overseer/overseer"
)

var allowList = []string{
	"/home/aco/code/Maelkum/overseer/cmd/example-server/example-server",
}

type actionFunc func(id string) (overseer.JobState, error)

func main() {

	var (
		flagCgroup string
		flagAllow  []string
		flagDeny   []string
	)

	pflag.StringVarP(&flagCgroup, "cgroup", "c", "overseer", "default cgroup to use for the overseer")
	pflag.StringArrayVarP(&flagAllow, "allow", "a", allowList, "list of binaries allowed to execute")
	pflag.StringArrayVarP(&flagDeny, "deny", "d", []string{}, "list of binaries not allowed to execute")

	pflag.Parse()

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	if flagCgroup == "" {
		log.Fatal().Msg("cgroup cannot be empty")
	}

	if !strings.HasPrefix(flagCgroup, "/") {
		flagCgroup = "/" + flagCgroup
	}

	limiter, err := limits.New(log, limits.DefaultMountpoint, flagCgroup)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create a limiter")
	}

	ov, err := overseer.New(
		log,
		limiter,
		overseer.WithAllowlist(flagAllow),
		overseer.WithDenylist(flagDeny),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create an overseer")
	}

	log.Info().Msg("starting job now")

	jobs := getJobs()

	jobHandles := make(map[string]*overseer.Handle)
	_ = jobHandles

	for i, job := range jobs {

		h, err := ov.Start(job)
		if err != nil {
			log.Fatal().Err(err).Int("index", i).Msg("could not start job")
		}

		log.Info().Int("index", i).Str("id", h.ID).Msg("job started")
	}

	log.Info().Msg("all jobs started")

	var ids []string
	for _, job := range jobs {
		ids = append(ids, job.ID)
	}

	// Control loop.
	for {
		jobID, err := selectOneOf("Select job", ids)
		if err != nil {
			log.Error().Err(err).Msg("job selection failed")
			break
		}

		fmt.Printf("Job chosen: %v\n", jobID)

		action, err := selectOneOf(
			"Choose job option",
			[]string{
				"stats",
				"wait",
				"kill",
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("action selection failed")
			break
		}

		var fn actionFunc

		switch action {
		case "stats":
			fn = ov.Stats

		case "wait":
			fn = ov.Wait

		case "kill":
			fn = ov.Kill
		}

		state, err := fn(jobID)
		if err != nil {
			log.Error().Err(err).Str("action", action).Msg("could not execute job action")
			continue
		}

		output, _ := json.MarshalIndent(state, "", "\t")
		fmt.Printf("%s\n", output)
	}

	for _, jobID := range ids {
		state, err := ov.Kill(jobID)
		if err != nil {
			log.Error().Err(err)
		}

		log.Info().Str("job", jobID).Interface("state", state).Msg("killed job")
	}

	log.Info().Msg("all done")
}

func selectOneOf(message string, opts []string) (string, error) {

	prompt := promptui.Select{
		Label: message,
		Items: opts,
	}

	for {
		_, selection, err := prompt.Run()
		if err != nil && errors.Is(err, promptui.ErrInterrupt) {
			return "", errors.New("interrupt")
		}
		if err != nil {
			log.Warn().Err(err).Msg("invalid selection")
			continue
		}

		return selection, nil
	}
}

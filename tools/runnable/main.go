package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {

	var (
		flagStdout         string
		flagStderr         string
		flagExitCode       int
		flagDuration       time.Duration
		flagBusywork       bool
		flagAllocateMemory int64
	)

	pflag.StringVar(&flagStdout, "stdout", "", "output to return on STDOUT")
	pflag.StringVar(&flagStderr, "stderr", "", "output to return on STDERR")
	pflag.IntVar(&flagExitCode, "exit-code", 0, "exit code for the app")

	pflag.DurationVar(&flagDuration, "duration", 0, "how long should the app run")
	pflag.Int64Var(&flagAllocateMemory, "allocate-memory", 0, "minimum memory the app should allocate")
	pflag.BoolVar(&flagBusywork, "busywork", false, "should the app do busywork to spend CPU cycles while running")

	pflag.Parse()

	defer func() {
		exitCode = flagExitCode
	}()

	// TODO: Check - do we try write if we have an empty string? Does that actually try anything?

	_, err := os.Stdout.Write([]byte(flagStdout))
	if err != nil {
		log.Printf("could not write to stdout: %s", err)
		return 1
	}

	_, err = os.Stderr.Write([]byte(flagStderr))
	// XXX: eems dubious this one would succeed.
	if err != nil {
		log.Printf("could not write to stderr: %s", err)
		return 1
	}

	if flagAllocateMemory > 0 {
		// TODO: See if this gets optimized away if not used in the loop.
		// If this does not get optimized away - this isn't necessary
		mem := make([]byte, flagAllocateMemory)
		for _, b := range mem {
			_ = b
		}
	}

	ch := make(chan struct{})

	go func() {

		if !flagBusywork {
			time.Sleep(flagDuration)
			close(ch)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), flagDuration)
		defer cancel()

		var x int

		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				x++
			}
		}
	}()

	<-ch

	return 0
}

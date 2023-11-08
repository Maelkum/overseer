package overseer

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"nhooyr.io/websocket"

	"github.com/Maelkum/overseer/limits"
)

type Handle struct {
	*sync.Mutex
	ID     string
	source Job

	stdout       *bytes.Buffer
	outputStream *websocket.Conn
	stderr       *bytes.Buffer
	errorStream  *websocket.Conn

	start     time.Time
	lastCheck time.Time

	cmd *exec.Cmd
}

func (o *Overseer) Start(job Job) (*Handle, error) {

	err := o.prepareJob(job)
	if err != nil {
		return nil, fmt.Errorf("could not prepare job: %w", err)
	}

	err = o.checkPrerequisites(job)
	if err != nil {
		return nil, fmt.Errorf("prerequisites not met: %w", err)
	}

	h, err := o.startJob(job)
	if err != nil {
		return nil, fmt.Errorf("could not start job: %w", err)
	}

	o.Lock()
	defer o.Unlock()

	o.jobs[job.ID] = h

	return h, nil
}

func (o *Overseer) startJob(job Job) (*Handle, error) {

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd, err := o.createCmd(job)
	if err != nil {
		return nil, fmt.Errorf("could not create command: %w", err)
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = job.Stdin

	o.log.Info().Str("cmd", cmd.String()).Msg("observer created command")

	handle := Handle{
		Mutex:  &sync.Mutex{},
		ID:     job.ID,
		source: job,

		stdout: &stdout,
		stderr: &stderr,

		cmd: cmd,
	}

	// Create an output stream if needed.
	if job.OutputStream != "" {
		// Continue even if stdout stream cannot be established.
		outputStream, err := wsConnect(job.OutputStream)
		if err != nil {
			o.log.Error().Err(err).Str("job", job.ID).Msg("could not establish output stream")
		} else {

			ws := wsWriter{
				conn: outputStream,
				log:  o.log.With().Str("job", job.ID).Logger(),
			}
			handle.outputStream = outputStream

			// Use both writers - both keep locally and stream data.
			// Websocket writer will never return errors as it's less important.
			cmd.Stdout = io.MultiWriter(&stdout, &ws)
		}
	}

	// Create an error stream too, if needed.
	if job.ErrorStream != "" {
		errorStream, err := wsConnect(job.ErrorStream)
		if err != nil {
			o.log.Error().Err(err).Str("job", job.ID).Msg("could not establish error stream")
		} else {

			ws := wsWriter{
				conn: errorStream,
				log:  o.log.With().Str("job", job.ID).Logger(),
			}
			handle.errorStream = errorStream

			cmd.Stderr = io.MultiWriter(&stderr, &ws)
		}
	}

	start := time.Now()
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start job: %w", err)
	}

	handle.start = start
	handle.lastCheck = start

	return &handle, nil
}

func (o *Overseer) prepareJob(job Job) error {

	workdir := o.workdir(job.ID)
	err := o.cfg.FS.MkdirAll(workdir, defaultFSPermissions)
	if err != nil {
		return fmt.Errorf("could not create work directory for request: %w", err)
	}

	return nil
}

func (o *Overseer) createCmd(job Job) (*exec.Cmd, error) {

	workdir := o.workdir(job.ID)

	cmd := exec.Command(job.Exec.Path, job.Exec.Args...)
	cmd.Dir = workdir
	cmd.Env = append(cmd.Env, job.Exec.Env...)

	if job.Limits != nil {

		opts := getLimitOpts(*job.Limits)
		err := o.limiter.CreateGroup(job.ID, opts...)
		if err != nil {
			return nil, fmt.Errorf("could not create limit group for job: %w", err)
		}

		fd, err := o.limiter.GetHandle(job.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get limit group handle: %w", err)
		}

		// NOTE: Setting child limits - https://man7.org/linux/man-pages/man2/clone3.2.html
		// Relevant:
		//	This file descriptor can be obtained by opening a cgroup v2 directory using either the O_RDONLY or the O_PATH flag.
		procAttr := syscall.SysProcAttr{
			UseCgroupFD: true,
			CgroupFD:    int(fd),
		}
		cmd.SysProcAttr = &procAttr
	}

	if o.cfg.FilesystemIsolation {

		if cmd.SysProcAttr != nil {
			cmd.SysProcAttr.Chroot = workdir
		} else {
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Chroot: workdir,
			}
		}
	}

	return cmd, nil
}

func (o *Overseer) workdir(id string) string {
	return filepath.Join(o.cfg.Workdir, id)
}

func getLimitOpts(jobLimits JobLimits) []limits.LimitOption {

	var opts []limits.LimitOption
	if jobLimits.CPUPercentage > 0 {
		opts = append(opts, limits.WithCPUPercentage(jobLimits.CPUPercentage))
	}

	if jobLimits.MemoryLimitKB > 0 {
		opts = append(opts, limits.WithMemoryKB(int64(jobLimits.MemoryLimitKB)))
	}

	if jobLimits.NoExec {
		opts = append(opts, limits.WithProcLimit(1))
	}

	return opts
}

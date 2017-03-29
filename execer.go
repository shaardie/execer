package execer

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"sync"
)

// Execer represent the Command to run.
type Execer struct {
	cmd      []string
	err      error
	finished bool
	mutex    *sync.Mutex
	started  bool
	stdout   string
	stderr   string
}

// Status represent the (temporarily) status of a `Execer`.
//
// `Error` stores the (last) error.
// `Finished` indicates if the command still running.
// `Started` indicates if the command is already started.
// `Stdout` contains the commands stdout until the time the `Status` was created.
// `Stderr` contains the commands stderr until the time the `Status` was created.
type Status struct {
	Error    error
	Finished bool
	Started  bool
	Stdout   string
	Stderr   string
}

// Create a new `Execer`. `len(cmd) == 0` is not allowed.
func Init(cmd []string) (Execer, error) {
	e := Execer{}
	if len(cmd) == 0 {
		return e, errors.New("command can not be empty")
	}
	e.mutex = &sync.Mutex{}
	e.cmd = cmd
	return e, nil
}

// Start the command.
// Returns an error if the command is not successfully started.
func (e *Execer) Start() (err error) {
	e.mutex.Lock()
	if e.started {
		err = errors.New("Unable to start twice")
		return
	}
	e.started = true
	e.mutex.Unlock()
	go e.runCmd()
	return
}

// Get the current status of the command run by `Execer`.
func (e Execer) Status() Status {
	return Status{
		e.err,
		e.finished,
		e.started,
		e.stdout,
		e.stderr}
}

func readTo(reader io.Reader, result *string, mutex *sync.Mutex) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		mutex.Lock()
		*result = *result + scanner.Text() + "\n"
		mutex.Unlock()
	}
}

func (e *Execer) runCmd() {
	runner := exec.Command(e.cmd[0], e.cmd[1:]...)
	defer func() { e.finished = true }()

	if reader, err := runner.StdoutPipe(); err == nil {
		go readTo(reader, &e.stdout, e.mutex)
	} else {
		e.err = err
		return
	}

	if reader, err := runner.StderrPipe(); err == nil {
		go readTo(reader, &e.stderr, e.mutex)
	} else {
		e.err = err
		return
	}

	if err := runner.Start(); err != nil {
		e.err = err
		return
	}

	if err := runner.Wait(); err != nil {
		e.err = err
	}
}

package execer

import (
	"bufio"
	"errors"
	"os/exec"
	"sync"
)

type Execer struct {
	cmd      []string
	err      error
	finished bool
	mutex    *sync.Mutex
	started  bool
	stdout   string
	stderr   string
}

type Status struct {
	Error    error
	Finished bool
	Started  bool
	Stdout   string
	Stderr   string
}

func Init(cmd []string) (Execer, error) {
	e := Execer{}
	if len(cmd) == 0 {
		return e, errors.New("command can not be empty")
	}
	e.mutex = &sync.Mutex{}
	e.cmd = cmd
	return e, nil
}

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

func (e Execer) Status() Status {
	return Status{
		e.err,
		e.finished,
		e.started,
		e.stdout,
		e.stderr}
}

func readTo(scanner *bufio.Scanner, result *string, mutex *sync.Mutex) {
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
		go readTo(bufio.NewScanner(reader), &e.stdout, e.mutex)
	} else {
		e.err = err
		return
	}

	if reader, err := runner.StderrPipe(); err == nil {
		go readTo(bufio.NewScanner(reader), &e.stderr, e.mutex)
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

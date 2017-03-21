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

func (e *Execer) runCmd() {
	runner := exec.Command(e.cmd[0], e.cmd[1:]...)
	defer func() { e.finished = true }()
	reader, err := runner.StdoutPipe()
	if err != nil {
		e.err = err
		return
	}

	scanner := bufio.NewScanner(reader)
	go func() {
		for scanner.Scan() {
			e.mutex.Lock()
			e.stdout = e.stdout + scanner.Text() + "\n"
			e.mutex.Unlock()
		}
	}()

	reader, err = runner.StderrPipe()
	if err != nil {
		e.err = err
		return
	}

	scanner2 := bufio.NewScanner(reader)
	go func() {
		for scanner2.Scan() {
			e.mutex.Lock()
			e.stderr = e.stderr + scanner2.Text() + "\n"
			e.mutex.Unlock()
		}
	}()

	err = runner.Start()
	if err != nil {
		e.err = err
		return
	}

	err = runner.Wait()
	if err != nil {
		e.err = err
	}

}

package execer

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestInitEmpty(t *testing.T) {
	_, err := Init([]string{})
	if err == nil {
		t.Error("empty cmd allowed")
	}
	t.Logf("Logged error: '%v'", err)
}

func TestInit(t *testing.T) {
	_, err := Init([]string{"true"})
	if err != nil {
		t.Error(err)
	}
}

func TestStart(t *testing.T) {
	e, err := Init([]string{"true"})
	if err != nil {
		t.Fatal(e.err)
	}
	if e.err != nil {
		t.Fatal(e.err)
	}
	err = e.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = e.Start()
	if err == nil {
		t.Fatal("Able to start twice")
	}
}

func TestStatus(t *testing.T) {
	e := Execer{}
	e.started = true
	if !e.Status().Started {
		t.Error("Wrong Started")
	}
	e.finished = true
	if !e.Status().Finished {
		t.Error("Wrong Finished")
	}
	e.stdout = "out"
	if e.Status().Stdout != "out" {
		t.Error("Wrong Stdout")
	}
	e.stderr = "err"
	if e.Status().Stderr != "err" {
		t.Error("Wrong Stdout")
	}
	err := errors.New("Fake Error")
	e.err = err
	if e.Status().Error != err {
		t.Error("Wrong Error")
	}
}

func TestReadTo(t *testing.T) {
	text := "This is\nen test\nfür dich!"
	text2 := text + "\n"
	var s string
	readTo(strings.NewReader(text), &s, &sync.Mutex{})
	if s != text2 {
		t.Errorf("Wrong text. Got '%v', but should be '%v'", s, text2)
	}
	// reset string
	s = ""
	readTo(strings.NewReader(text2), &s, &sync.Mutex{})
	if s != text2 {
		t.Errorf("Wrong text. Got '%v', but should be '%v'", s, text2)
	}
}

func TestRunCmd(t *testing.T) {
	s := "This is a\n\nmulti line text!\n"
	e, err := Init([]string{"./testcmd", "--stderr", "--stdout", "--text", s})
	e.runCmd()
	if e.err != nil {
		t.Error(err)
	}
	if e.stdout != s+"\n" {
		t.Errorf("Wrong stdout. Got '%v'", e.stdout)
	}
	if e.stderr != s+"\n" {
		t.Errorf("Wrong stderr. Got '%v'", e.stderr)
	}
}

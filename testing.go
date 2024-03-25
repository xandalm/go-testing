package testing

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type AvailabilityChecker interface {
	// Should return error if unable to ping (didn't pong)
	Ping() error
}

type HTTPServerChecker struct {
	BaseURL string
	Cli     *http.Client
}

func (c *HTTPServerChecker) Ping() error {
	if _, err := c.Cli.Get(c.BaseURL); err != nil {
		return err
	}
	return nil
}

type ServerLauncher struct {
	ctx  context.Context
	wd   string
	name string
	c    AvailabilityChecker
	cmd  *exec.Cmd
}

func NewServerLauncher(ctx context.Context, wd, filename string, checker AvailabilityChecker) *ServerLauncher {
	if ctx == nil {
		panic("testing: nil context")
	}
	if checker == nil {
		panic("testing: nil server availability checker")
	}
	if filename == "" {
		panic("testing: empty file name")
	}
	if !strings.HasSuffix(filename, ".go") {
		panic("testing: must be a go file (.go)")
	}
	return &ServerLauncher{ctx, wd, strings.TrimSuffix(filename, ".go"), checker, nil}
}

func ping(c AvailabilityChecker) chan bool {
	ch := make(chan bool, 1)
	go func() {
		ch <- c.Ping() == nil
	}()
	return ch
}

func (s *ServerLauncher) build() error {
	cmd := exec.CommandContext(s.ctx, "go", "build", s.name+".go")
	cmd.Dir = s.wd
	return cmd.Run()
}

func (s *ServerLauncher) wait() chan struct{} {
	ch := make(chan struct{})
	for {
		select {
		case res := <-ping(s.c):
			if res {
				close(ch)
				return ch
			}
		case <-s.ctx.Done():
			close(ch)
			return ch
		}
	}
}

func (s *ServerLauncher) clean() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(s.ctx, "cmd", "/C", fmt.Sprintf("del %s.exe", s.name))
	} else {
		cmd = exec.CommandContext(s.ctx, "rm", "./"+s.name)
	}
	cmd.Dir = s.wd
	return cmd.Run()
}

func (s *ServerLauncher) EndAndClean() error {
	if err := s.cmd.Cancel(); err != nil {
		return fmt.Errorf("testing: cannot end server and clean")
	}
	if err := s.clean(); err != nil {
		return fmt.Errorf("testing: cannot clean")
	}
	return nil
}

func (s *ServerLauncher) StartAndWait(waitFor time.Duration) error {
	if err := s.build(); err != nil {
		return fmt.Errorf("testing: cannot build and start server, %v", err)
	}

	s.cmd = exec.CommandContext(s.ctx, "./"+s.name)
	s.cmd.Dir = s.wd

	go func() {
		if err := s.cmd.Run(); err != nil {
			s.clean()
		}
	}()

	select {
	case <-s.wait():
		return nil
	case <-time.After(waitFor):
		s.EndAndClean()
		return fmt.Errorf("testing: cannot start server")
	}
}

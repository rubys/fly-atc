package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type UpstreamProcess struct {
	Started chan struct{}
	cmd     *exec.Cmd
}

func NewUpstreamProcess(name string, arg ...string) *UpstreamProcess {
	cmd := exec.Command(name, arg...)
	cmd.Env = os.Environ()

	return &UpstreamProcess{
		Started: make(chan struct{}, 1),
		cmd:     cmd,
	}
}

func (p *UpstreamProcess) setEnvironment(name string, value string) {
	p.cmd.Env = append(p.cmd.Env, fmt.Sprintf("%s=%s", name, value))
}

func (p *UpstreamProcess) Start() error {
	p.cmd.Stdin = os.Stdin
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stderr

	err := p.cmd.Start()
	if err != nil {
		return err
	}

	p.Started <- struct{}{}

	go p.handleSignals()

	return nil
}

func (p *UpstreamProcess) Stop() (int, error) {
	err := p.cmd.Wait()

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), nil
	}

	return 0, err
}

func (p *UpstreamProcess) Signal(sig os.Signal) error {
	return p.cmd.Process.Signal(sig)
}

func (p *UpstreamProcess) handleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	sig := <-ch
	slog.Info("Relaying signal to upstream process", "signal", sig.String())
	p.Signal(sig)
}

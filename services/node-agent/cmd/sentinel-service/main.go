//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sys/windows/svc"
)

const serviceName = "sentinel-agent"

type service struct {
	process *os.Process
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	changes <- svc.Status{State: svc.StartPending}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get executable path: %v\n", err)
		changes <- svc.Status{State: svc.Stopped}
		return false, 1
	}

	dir := filepath.Dir(exePath)
	agentPath := filepath.Join(dir, "sentinel-agent.exe")

	cmd := exec.Command(agentPath, "run")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start sentinel-agent: %v\n", err)
		changes <- svc.Status{State: svc.Stopped}
		return false, 1
	}

	s.process = cmd.Process
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				if s.process != nil {
					_ = s.process.Signal(syscall.SIGTERM)
					done := make(chan struct{})
					go func() {
						_, _ = s.process.Wait()
						close(done)
					}()
					select {
					case <-done:
					case <-time.After(10 * time.Second):
						if s.process != nil {
							_ = s.process.Kill()
						}
					}
				}
				changes <- svc.Status{State: svc.Stopped}
				return false, 0
			}
		}
	}
}

func main() {
	svc.IsWindowsService()
	interactive, err := svc.IsAnInteractiveSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to detect session: %v\n", err)
		os.Exit(1)
	}

	if interactive {
		fmt.Fprintln(os.Stderr, "running in interactive mode, use install to run as service")
		os.Exit(0)
	}

	svc.Run(serviceName, &service{})
}

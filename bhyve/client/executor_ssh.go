package client

import (
	"errors"
	"fmt"
	"os/exec"
)

type SSHExecutor struct {
	Rsh    string
	Host   string
	Port   int
	User   string
	Params []string
}

var _ Executor = (*SSHExecutor)(nil)

type RemoteExecutionError struct {
	err error
	Cmd *exec.Cmd
}

func (e RemoteExecutionError) Error() string {
	return e.err.Error()
}

func (e SSHExecutor) constructCommand(cmd string, args ...string) (*exec.Cmd, error) {
	rsh := e.Rsh
	if rsh == "" {
		rsh = "ssh"
	}

	host := e.Host
	if host == "" {
		return nil, errors.New("host must be specified")
	}

	port := e.Port
	if port == 0 {
		port = 22
	}

	user := e.User

	finalParams := []string{host, "-p", fmt.Sprint(port)}
	if user != "" {
		finalParams = append(finalParams, "-l", user)
	}

	finalParams = append(finalParams, e.Params...)
	finalParams = append(finalParams, "--", cmd)
	finalParams = append(finalParams, args...)

	// TODO: add params to set up connection reuse

	return exec.Command(rsh, finalParams...), nil
}

func (e SSHExecutor) RunCmd(cmd string, args ...string) (string, error) {
	execCmd, err := e.constructCommand(cmd, args...)
	if err != nil {
		return "", err
	}

	if out, err := execCmd.Output(); err != nil {
		stderr := ""
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = string(ee.Stderr)
		}
		return "", CommandExecutionError{
			Wrapped:    err,
			ReturnCode: execCmd.ProcessState.ExitCode(),

			Cmd:    execCmd.Args,
			Stderr: stderr,
		}
	} else {
		return string(out), nil
	}
}

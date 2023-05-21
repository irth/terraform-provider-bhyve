package client

type CommandExecutionError struct {
	// Wrapped contains the underlying error, if any.
	Wrapped error

	// Cmd contains the full command that was executed.
	Cmd []string

	// ReturnCode contains the return code of the process.
	ReturnCode int

	/// Stderr contains the stderr output of the process.
	Stderr string
}

func (e CommandExecutionError) Error() string {
	return e.Wrapped.Error()
}

type Executor interface {
	RunCmd(cmd string, args ...string) (string, error)
}

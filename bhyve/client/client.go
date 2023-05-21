package client

import "errors"

var ErrInvalidOutput = errors.New("vm command returned invalid output")

type Client struct {
	Executor
}

package client

import "errors"

var ErrInvalidOutput = errors.New("vm command returned invalid output")
var ErrNotFound = errors.New("resource not found")

type Client struct {
	Executor
}

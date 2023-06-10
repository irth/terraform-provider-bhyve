package client

import (
	"errors"
	"regexp"
	"strings"
)

type Sysrc struct {
	Executor Executor
	File     string
}

func (s Sysrc) run(args ...string) (string, error) {
	var finalArgs = []string{"-q"}
	if s.File != "" {
		finalArgs = []string{"-f", s.File}
	}
	finalArgs = append(finalArgs, args...)
	return s.Executor.RunCmd("sysrc", finalArgs...)
}

var validKey = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

var ErrInvalidKey = errors.New("key contains invalid characters")

// validateKeys checks if the key only contains valid characters.
func (s Sysrc) validateKey(key string) error {
	if !validKey.MatchString(key) {
		return ErrInvalidKey
	}
	return nil
}

var ErrUnknownVariable = errors.New("unknown variable")

func (s Sysrc) Get(key string) (string, error) {
	if err := s.validateKey(key); err != nil {
		return "", err
	}
	v, err := s.run("-n", key)
	if err != nil {
		if err, ok := err.(CommandExecutionError); ok {
			if err.ReturnCode == 1 {
				// I'm assuming 1 means unknown, because that's how it behaves
				// for me, but I haven't found this in the manpage
				// TODO: check if this is correct
				return "", ErrUnknownVariable
			}
		}
		return "", err
	}

	return strings.TrimSpace(v), nil
}

func (s Sysrc) GetDefault(key, defaultValue string) (string, error) {
	v, err := s.Get(key)
	if err != nil {
		if err == ErrUnknownVariable {
			err = nil
		}
		return defaultValue, err
	}
	return v, nil
}

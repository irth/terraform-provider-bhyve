package client

import (
	"errors"
	"strings"

	"github.com/samber/lo"
)

type Switch struct {
	Name    string
	Address string
}

type Switches []Switch

var ErrInvalidOutput = errors.New("vm command returned invalid output")

func (s *Switches) LoadFromSystem(executor Executor) error {
	out, err := executor.RunCmd("vm", "switch", "list")
	if err != nil {
		return err
	}

	lines := strings.Split(out, "\n")
	header, lines := lines[0], lines[1:]
	headerFields := strings.Fields(header)

	nameIdx := lo.IndexOf(headerFields, "NAME")
	addrIdx := lo.IndexOf(headerFields, "ADDRESS")

	if nameIdx == -1 || addrIdx == -1 {
		return ErrInvalidOutput
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != len(headerFields) {
			continue
		}
		name, addr := fields[nameIdx], fields[addrIdx]
		if addr == "-" {
			addr = ""
		}
		*s = append(*s, Switch{Name: name, Address: addr})
	}
	return nil
}

func (s Switches) AsMap() map[string]Switch {
	m := make(map[string]Switch, len(s))
	for _, sw := range s {
		m[sw.Name] = sw
	}
	return m
}

package client

import (
	"strings"

	"github.com/samber/lo"
)

type Switch struct {
	Name    string
	Address string
}

func (c *Client) SwitchList() (map[string]Switch, error) {
	// TODO: use ctx
	switches := make(map[string]Switch)

	out, err := c.RunCmd("vm", "switch", "list")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")
	header, lines := lines[0], lines[1:]
	headerFields := strings.Fields(header)

	nameIdx := lo.IndexOf(headerFields, "NAME")
	addrIdx := lo.IndexOf(headerFields, "ADDRESS")

	if nameIdx == -1 || addrIdx == -1 {
		return nil, ErrInvalidOutput
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
		switches[name] = Switch{Name: name, Address: addr}
	}
	return switches, nil
}

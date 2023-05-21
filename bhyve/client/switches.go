package client

import (
	"strings"

	"github.com/samber/lo"
)

type Switch struct {
	Name    string
	Address string

	client *Client
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
		switches[name] = Switch{Name: name, Address: addr, client: c}
	}
	return switches, nil
}

func (c *Client) SwitchCreate(sw *Switch) error {
	sw.client = c
	// TODO: validate address is cidr

	params := []string{"switch", "create", sw.Name}
	if sw.Address != "" {
		params = append(params, "-a", sw.Address)
	}

	_, err := c.RunCmd("vm", params...)
	if err != nil {
		return err
	}

	if sw.Address != "" {
		// for some reason, the `-a` flag doesn't always work
		err = c.SwitchAddress(sw.Name, sw.Address)
	}

	return err
}

func (c *Client) SwitchDestroy(name string) error {
	params := []string{"switch", "destroy", name}
	_, err := c.RunCmd("vm", params...)
	return err
}

func (c *Client) SwitchInfo(name string) (*Switch, error) {
	switches, err := c.SwitchList()
	if err != nil {
		return nil, err
	}

	sw, ok := switches[name]
	if !ok {
		return nil, ErrNotFound
	}

	return &sw, nil
}

func (c *Client) SwitchAddress(name string, addr string) error {
	if addr == "" {
		addr = "none"
	}
	_, err := c.RunCmd("vm", "switch", "address", name, addr)
	return err
}

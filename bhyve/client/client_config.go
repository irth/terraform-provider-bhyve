package client

import (
	"fmt"
	"path"
	"strings"
)

type Config struct {
	BhyveEnabled bool
	VMEnabled    bool
	VMDir        string
}

func (c *Client) Config() (Config, error) {
	cfg := Config{}
	sysrc := Sysrc{Executor: c}

	var err error
	cfg.BhyveEnabled, err = sysrc.GetBoolDefault("vmm_load", false)
	if err != nil {
		return cfg, err
	}

	cfg.VMEnabled, err = sysrc.GetBoolDefault("vm_enable", false)
	if err != nil {
		return cfg, err
	}

	cfg.VMDir, err = sysrc.GetDefault("vm_dir", "")
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c *Config) Path(base string, elem ...string) string {
	// TODO(wkwolek): move to using VM dir from provider config, to avoid figuring out ZFS lmao
	vmDir := strings.TrimPrefix(c.VMDir, "zfs:")
	if vmDir[0] != '/' {
		vmDir = fmt.Sprintf("/%s", vmDir)
	}
	fullPath := []string{vmDir, base}
	fullPath = append(fullPath, elem...)
	return path.Join(fullPath...)
}

func (c *Config) ISOPath(elem ...string) string {
	return c.Path(".iso", elem...)
}

func (c *Config) IMGPath(elem ...string) string {
	return c.Path(".img", elem...)
}

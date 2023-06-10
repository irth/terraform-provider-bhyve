package client

import (
	"fmt"
	"path"
	"strings"
)

type Config struct {
	BhyveEnabled bool   `sysrc:"vmm_load"`
	VMEnabled    bool   `sysrc:"vm_enable"`
	VMDir        string `sysrc:"vm_dir"`
}

func (c *Client) Config() (Config, error) {
	cfg := Config{}
	sysrc := Sysrc{Executor: c}
	sysrc.Unmarshal(&cfg)

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

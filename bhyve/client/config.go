package client

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

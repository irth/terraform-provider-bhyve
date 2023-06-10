package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestSysrcExecutor struct {
	Vals map[string]string
}

func (e TestSysrcExecutor) RunCmd(cmd string, args ...string) (string, error) {
	if cmd != "sysrc" {
		return "", fmt.Errorf("unexpected command: %s", cmd)
	}

	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			continue
		}

		val, ok := e.Vals[arg]
		if ok {
			return val, nil
		}
	}

	return "", ErrUnknownVariable
}

func TestUnmarshalSysrc(t *testing.T) {
	var target struct {
		VmEnable bool   `sysrc:"vm_enable"`
		VmDir    string `sysrc:"vm_dir"`
	}

	sysrc := Sysrc{
		Executor: TestSysrcExecutor{
			Vals: map[string]string{
				"vm_enable": "YES",
				"vm_dir":    "/vm",
			},
		},
		File: "/etc/rc.conf",
	}
	err := sysrc.Unmarshal(&target)
	assert.NoError(t, err)
	assert.Equal(t, true, target.VmEnable)
	assert.Equal(t, "/vm", target.VmDir)

	sysrc = Sysrc{
		Executor: TestSysrcExecutor{
			Vals: map[string]string{
				"vm_enable": "NO",
				"vm_dir":    "/test",
			},
		},
		File: "/etc/rc.conf",
	}

	err = sysrc.Unmarshal(&target)
	assert.NoError(t, err)
	assert.Equal(t, false, target.VmEnable)
	assert.Equal(t, "/test", target.VmDir)

}

func TestUnmarshalSysrcMissing(t *testing.T) {
	var target struct {
		VmEnable bool   `sysrc:"vm_enable"`
		VmDir    string `sysrc:"vm_dir"`
	}

	sysrc := Sysrc{
		Executor: TestSysrcExecutor{
			Vals: map[string]string{},
		},
		File: "/etc/rc.conf",
	}
	err := sysrc.Unmarshal(&target)
	assert.NoError(t, err)
	assert.Equal(t, false, target.VmEnable)
	assert.Equal(t, "", target.VmDir)
}

func TestUnmarshalSysrcMissingRequired(t *testing.T) {
	var target struct {
		VmEnable bool   `sysrc:"vm_enable,required"`
		VmDir    string `sysrc:"vm_dir"`
	}

	sysrc := Sysrc{
		Executor: TestSysrcExecutor{
			Vals: map[string]string{},
		},
		File: "/etc/rc.conf",
	}
	err := sysrc.Unmarshal(&target)
	assert.Error(t, err)
}

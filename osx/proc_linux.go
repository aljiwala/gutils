package osx

import (
	"os/exec"
	"syscall"
)

func procAttrSetGroup(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true, // to be able to kill all children, too
		Pdeathsig: syscall.SIGKILL,
	}
}

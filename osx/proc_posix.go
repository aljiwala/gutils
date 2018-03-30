package osx

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func isGroupLeader(c *exec.Cmd) bool {
	return c.SysProcAttr != nil && c.SysProcAttr.Setpgid
}

// Pkill kills the process with the given pid, or just -INT if interrupt is true.
func Pkill(pid int, signal os.Signal) error {
	signum := signal.(syscall.Signal)

	var err error
	defer func() {
		if r := recover(); r == nil && err == nil {
			return
		}
		err = exec.Command("pkill", "-"+strconv.Itoa(int(signum)),
			"-P", strconv.Itoa(pid)).Run()
	}()
	err = syscall.Kill(pid, signum)
	return err
}

// GroupKill kills the process group lead by the given pid
func GroupKill(pid int, signal os.Signal) error {
	return Pkill(-pid, signal)
}

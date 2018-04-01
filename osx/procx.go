package osx

import (
	"fmt"
	"os"
	"os/exec"
)

// KillWithChildren kills the process
// and tries to kill its all children (process group)
func KillWithChildren(p *os.Process, interrupt bool) (err error) {
	if p == nil {
		return
	}
	fmt.Println("msg", "killWithChildren", "pid", p.Pid, "interrupt", interrupt)
	defer func() {
		if r := recover(); r != nil {
			Log("msg", "PANIC in kill", "process", p, "error", r)
		}
	}()
	defer p.Release()

	if p.Pid == 0 {
		return nil
	}
	if interrupt {
		defer p.Signal(os.Interrupt)
		return Pkill(p.Pid, os.Interrupt)
	}
	defer p.Kill()

	return Pkill(p.Pid, os.Kill)
}

func groupKill(p *os.Process, interrupt bool) error {
	if p == nil {
		return nil
	}
	fmt.Println("msg", "groupKill", "pid", p.Pid)
	defer recover()

	if interrupt {
		defer p.Signal(os.Interrupt)
		return GroupKill(p.Pid, os.Interrupt)
	}
	defer p.Kill()

	return GroupKill(p.Pid, os.Kill)
}

func familyKill(cmd *exec.Cmd, interrupt bool) error {
	if cmd.SysProcAttr != nil && isGroupLeader(cmd) {
		return groupKill(cmd.Process, interrupt)
	}
	return KillWithChildren(cmd.Process, interrupt)
}

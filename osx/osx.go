package osx

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	lastmod, lastcheck time.Time
	groups             map[int]string
	groupsMu           sync.RWMutex
)

const groupFile = "/etc/group"

// IntTimeout is the duration to wait before Kill after Int
var IntTimeout = 3 * time.Second

// Log is discarded by default
var Log = func(keyvals ...interface{}) error { return nil }

// ErrTimedOut is an error for child timeout
var ErrTimedOut = errors.New("child timed out")

type gCmd struct {
	*exec.Cmd
	done chan error
}

func (c *gCmd) Start() error {
	if err := c.Cmd.Start(); err != nil {
		return err
	}
	c.done = make(chan error, 1)
	go func() { c.done <- c.Cmd.Wait() }()
	return nil
}

// RunWithTimeout runs cmd, and kills the child on timeout
func RunWithTimeout(timeoutSeconds int, cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		procAttrSetGroup(cmd)
	}

	gcmd := &gCmd{Cmd: cmd}
	if err := gcmd.Start(); err != nil {
		return err
	}
	if timeoutSeconds <= 0 {
		return <-gcmd.done
	}

	select {
	case err := <-gcmd.done:
		return err
	case <-time.After(time.Second * time.Duration(timeoutSeconds)):
		Log("msg", "killing timed out", "pid", cmd.Process.Pid, "path", cmd.Path, "args", cmd.Args)
		if killErr := familyKill(gcmd.Cmd, true); killErr != nil {
			Log("msg", "interrupt", "pid", cmd.Process.Pid)
		}

		select {
		case <-gcmd.done:
		case <-time.After(IntTimeout):
			familyKill(gcmd.Cmd, false)
		}
	}

	return ErrTimedOut
}

// GroupName returns the name for the gid.
func GroupName(gid int) (string, error) {
	groupsMu.RLock()
	if groups == nil {
		groupsMu.RUnlock()
		groupsMu.Lock()
		defer groupsMu.Unlock()
		if groups != nil { // sy was faster
			name := groups[gid]
			return name, nil
		}
	} else {
		now := time.Now()
		if lastcheck.Add(1 * time.Second).After(now) { // fresh
			name := groups[gid]
			groupsMu.RUnlock()
			return name, nil
		}

		actcheck := lastcheck
		groupsMu.RUnlock()
		groupsMu.Lock()
		defer groupsMu.Unlock()
		if lastcheck != actcheck { // sy was faster
			return groups[gid], nil
		}

		fi, err := os.Stat(groupFile)
		if err != nil {
			return "", err
		}

		lastcheck = now
		if lastmod == fi.ModTime() { // no change
			return groups[gid], nil
		}
	}

	// need to reread
	if groups == nil {
		groups = make(map[int]string, 64)
	} else {
		for k := range groups {
			delete(groups, k)
		}
	}

	fh, err := os.Open(groupFile)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	fi, err := fh.Stat()
	if err != nil {
		return "", err
	}

	lastcheck = time.Now()
	lastmod = fi.ModTime()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 4)
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Printf("cannot parse %q as group id from line %q: %v", parts[2], scanner.Text(), err)
		}
		if old, ok := groups[id]; ok {
			log.Printf("double entry %d: %q and %q?", id, old, parts[0])
			continue
		}
		groups[id] = parts[0]
	}

	return groups[gid], nil
}

// IsInsideDocker returns true iff we are inside a docker cgroup.
func IsInsideDocker() bool {
	b, err := ioutil.ReadFile("/proc/self/cgroup")
	if err != nil {
		return false
	}
	return bytes.Contains(b, []byte(":/docker/")) || bytes.Contains(b, []byte(":/lxc/"))
}

package cmdx

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Log is discarded by default
var Log = func(...interface{}) error { return nil }

// Loffice executable name
var Loffice = "loffice"

// Timeout of the child, in seconds
var Timeout = 300

// pipeCommands should return a pipe that will be connected to the command's
// standard output when the command starts.
func pipeCommands(commands ...*exec.Cmd) ([]byte, error) {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		command.Start()
		commands[i+1].Stdin = out
	}

	final, err := commands[len(commands)-1].Output()
	if err != nil {
		return nil, err
	}

	return final, nil
}

// OutStr should run given command with provided arguments; returns printed
// output with error (if any).
func OutStr(name string, arg ...string) (outStr string, err error) {
	var out bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &out
	err = cmd.Run()
	outStr = out.String()
	return
}

// OutBytes should run given command with provided arguments; returns printed
// output (as bytes) with error (if any).
func OutBytes(name string, arg ...string) (byteContainer []byte, err error) {
	var out bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &out
	err = cmd.Run()
	byteContainer = out.Bytes()
	return
}

// RunWithTimeout should run cmd with given timeout duration.
func RunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (bool, error) {
	var err error
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout): // Timeout.
		if err = cmd.Process.Kill(); err != nil {
			log.Printf("failed to kill: %s, error: %s", cmd.Path, err) // ERROR
		}
		go func() {
			<-done // Allow `goroutine` to exit.
		}()
		log.Printf("process:%s killed", cmd.Path) // INFO
		return true, err
	case err = <-done:
		return false, err
	}
}

// Convert converts from srcFn to dstFn, into the given format.
// Convert from srcFn to dstFn files, with the given format.
// Either filenames can be empty or "-" which treated as stdin/stdout
func ConvertLoffice(srcFn, dstFn, format string) error {
	tempDir, err := ioutil.TempDir("", filepath.Base(srcFn))
	if err != nil {
		return fmt.Errorf("cannot create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	if srcFn == "-" || srcFn == "" {
		srcFn = filepath.Join(tempDir, "source")
		fh, cErr := os.Create(srcFn)
		if cErr != nil {
			return fmt.Errorf("error creating temp file %q: %s", srcFn, cErr)
		}
		if _, err = io.Copy(fh, os.Stdin); err != nil {
			fh.Close()
			return fmt.Errorf("error writing stdout to %q: %s", srcFn, err)
		}
		fh.Close()
	}

	c := exec.Command(Loffice, "--nolockcheck", "--norestore", "--headless",
		"--convert-to", format, "--outdir", tempDir, srcFn)
	c.Stderr = os.Stderr
	c.Stdout = c.Stderr

	Log("msg", "calling", "args", c.Args)
	if err = proc.RunWithTimeout(Timeout, c); err != nil {
		return fmt.Errorf("error running %q: %s", c.Args, err)
	}

	dh, err := os.Open(tempDir)
	if err != nil {
		return fmt.Errorf("error opening dest dir %q: %s", tempDir, err)
	}
	defer dh.Close()

	names, err := dh.Readdirnames(3)
	if err != nil {
		return fmt.Errorf("error listing %q: %s", tempDir, err)
	}
	if len(names) > 2 {
		return fmt.Errorf("too many files in %q: %q", tempDir, names)
	}

	var tfn string
	for _, fn := range names {
		if fn != "source" {
			tfn = filepath.Join(dh.Name(), fn)
			break
		}
	}

	src, err := os.Open(tfn)
	if err != nil {
		return fmt.Errorf("cannot open %q: %s", tfn, err)
	}
	defer src.Close()

	var dst = io.WriteCloser(os.Stdout)
	if !(dstFn == "-" || dstFn == "") {
		if dst, err = os.Create(dstFn); err != nil {
			return fmt.Errorf("cannot create dest file %q: %s", dstFn, err)
		}
	}
	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("error copying from %v to %v: %v", src, dst, err)
	}

	return nil
}

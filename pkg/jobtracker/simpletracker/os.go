package simpletracker

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/dgruber/drmaa2interface"
)

func currentEnv() map[string]string {
	env := make(map[string]string, len(os.Environ()))
	for _, e := range os.Environ() {
		env[e] = os.Getenv(e)
	}
	return env
}

func restoreEnv(env map[string]string) {
	for _, e := range os.Environ() {
		os.Unsetenv(e)
	}
	for key, value := range env {
		os.Setenv(key, value)
	}
}

// StartProcess creates a new process based on the JobTemplate.
// It returns the PID or 0 and an error if the process could be
// created. The given channel is used for communicating back
// when the job state changed.
func StartProcess(jobid string, task int, t drmaa2interface.JobTemplate, finishedJobChannel chan JobEvent) (int, error) {
	cmd := exec.Command(t.RemoteCommand, t.Args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if valid, err := validateJobTemplate(t); valid == false {
		return 0, err
	}

	waitForFiles := 0
	waitCh := make(chan bool, 3)

	var mtx sync.Mutex
	mtx.Lock()
	defer mtx.Unlock()

	if t.InputPath != "" {

		if t.InputPath == "/dev/stdin" {
			cmd.Stdin = os.Stdin
		} else {
			file, err := os.Open(t.InputPath)
			if err != nil {
				panic(err)
			}
			cmd.Stdin = file
		}

		/*
			stdin, err := cmd.StdinPipe()
			if err == nil {
				err = redirectIn(stdin, t.InputPath, waitCh)
				if err == nil {
					waitForFiles++
				} else {
					panic(err)
				}
			}
		*/

	}
	if t.OutputPath != "" {

		if t.OutputPath == "/dev/stdout" {
			cmd.Stdout = os.Stdout
		} else if t.OutputPath == "/dev/stderr" {
			cmd.Stdout = os.Stderr
		} else {
			outfile, err := os.Create(t.OutputPath)
			if err != nil {
				panic(err)
			}
			cmd.Stdout = outfile
		}

		/*
			stdout, err := cmd.StdoutPipe()
			if err == nil {
				err = redirectOut(stdout, t.OutputPath, waitCh)
				if err == nil {
					waitForFiles++
				} else {
					panic(err)
				}
			}
		*/

	}
	if t.ErrorPath != "" {

		if t.ErrorPath == "/dev/stdout" {
			cmd.Stderr = os.Stdout
		} else if t.ErrorPath == "/dev/stderr" {
			cmd.Stderr = os.Stderr
		} else {
			outfile, err := os.Create(t.ErrorPath)
			if err != nil {
				panic(err)
			}
			cmd.Stdout = outfile
		}

		/*
			stderr, err := cmd.StderrPipe()
			if err == nil {
				err = redirectOut(stderr, t.ErrorPath, waitCh)
				if err == nil {
					waitForFiles++
				}
			}
		*/
	}

	cmd.Env = os.Environ()
	for key, value := range t.JobEnvironment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("JOB_ID=%s", jobid))
	if task != 0 {
		cmd.Env = append(cmd.Env, fmt.Sprintf("TASK_ID=%d", task))
	}

	cmd.Dir = t.WorkingDirectory

	if t.ExtensionList != nil {
		if _, exists := t.ExtensionList["chroot"]; exists {
			cmd.SysProcAttr.Chroot = t.ExtensionList["chroot"]
		}
	}

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	host, _ := os.Hostname()
	startTime := time.Now()

	finishedJobChannel <- JobEvent{
		JobState: drmaa2interface.Running,
		JobID:    jobid,
		JobInfo: drmaa2interface.JobInfo{
			State:             drmaa2interface.Running,
			DispatchTime:      startTime,
			AllocatedMachines: []string{host},
		},
	}

	go TrackProcess(cmd, nil, jobid, startTime, finishedJobChannel, waitForFiles, waitCh)

	if cmd.Process == nil {
		return 0, errors.New("process is nil")
	}
	return cmd.Process.Pid, nil
}

func redirectOut(src io.ReadCloser, outfilename string, waitCh chan bool) error {
	var outfile *os.File
	if outfilename == "/dev/stdout" {
		outfile = os.Stdout
	} else if outfilename == "/dev/stderr" {
		outfile = os.Stderr
	} else {
		var err error
		outfile, err = os.Create(outfilename)
		if err != nil {
			return err
		}
	}
	go func() {
		_, err := io.Copy(outfile, src)
		if err != nil {
			fmt.Printf("Internal error while copying buffer: %s", err)
		}
		src.Close()
		outfile.Close()
		waitCh <- true
	}()
	return nil
}

func redirectIn(out io.WriteCloser, infilename string, waitCh chan bool) error {
	file, err := os.Open(infilename)
	if err != nil {
		return err
	}
	go func() {
		//buf := make([]byte, 1024)
		io.Copy(out, file)
		file.Close()
		// need to close stdin otherwise cmd might wait infinitely
		out.Close()
		waitCh <- true
	}()
	return nil
}

// KillPid terminates a process and all processes belonging
// to the process group.
func KillPid(pid int) error {
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return syscall.Kill(-pid, syscall.SIGKILL)
	}
	return syscall.Kill(-pgid, syscall.SIGKILL)
}

// SuspendPid stops a process group from its execution. Note
// that it sends SIGTSTP which can be caught by the application
// and hence could be ignored.
func SuspendPid(pid int) error {
	return syscall.Kill(-pid, syscall.SIGTSTP)
}

// ResumePid contiues to run a previously suspended process group.
func ResumePid(pid int) error {
	return syscall.Kill(-pid, syscall.SIGCONT)
}

// IsPidRunning returns true if the process is still alive.
func IsPidRunning(pid int) (bool, error) {
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return false, err
	}
	err = process.Signal(syscall.Signal(0))
	if err == nil {
		return true, nil
	}
	if errno, ok := err.(syscall.Errno); ok {
		switch errno {
		case syscall.EPERM:
			return true, nil
		}
	}
	return false, nil
}

package simpletracker

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/extension"
)

// TrackProcess supervises a running process and sends a notification when
// the process is finished. If the process was started from this process
// cmd is given otherwise when when re-attaching to an already existing
// process proc is given.
func TrackProcess(cmd *exec.Cmd, proc *os.Process, jobID string, startTime time.Time,
	finishedJobChannel chan JobEvent, waitForFiles int, waitCh chan bool) {

	var state *os.ProcessState
	var err error

	if cmd != nil {
		// wait for process and get state
		state, err = cmd.Process.Wait()
		// additionally wait for file descriptors to prevent blocking
		cmd.Wait()
	} else {
		state, err = proc.Wait()
	}

	// wait until all filedescriptors (stdout, stderr) of the
	// process are closed
	if cmd != nil {
		for waitForFiles > 0 {
			<-waitCh
			waitForFiles--
		}
	}

	if err != nil {
		ji := makeLocalJobInfo()
		ji.State = drmaa2interface.Failed
		finishedJobChannel <- JobEvent{
			JobState: drmaa2interface.Failed,
			JobID:    jobID,
			JobInfo:  ji,
		}
		return
	}

	ji := collectUsage(state, jobID, startTime)
	finishedJobChannel <- JobEvent{JobState: ji.State, JobID: jobID, JobInfo: ji}
}

func makeLocalJobInfo() drmaa2interface.JobInfo {
	host, _ := os.Hostname()
	return drmaa2interface.JobInfo{
		AllocatedMachines: []string{host},
		FinishTime:        time.Now(),
		SubmissionMachine: host,
		JobOwner:          fmt.Sprintf("%d", os.Getuid()),
	}
}

func collectUsage(state *os.ProcessState, jobid string, startTime time.Time) drmaa2interface.JobInfo {
	ji := makeLocalJobInfo()
	ji.State = drmaa2interface.Failed

	if status, ok := state.Sys().(syscall.WaitStatus); ok {
		ji.ExitStatus = status.ExitStatus()
		ji.TerminatingSignal = status.Signal().String()
	}

	if ji.ExtensionList == nil {
		ji.ExtensionList = make(map[string]string)
	}

	if usage, ok := state.SysUsage().(syscall.Rusage); ok {
		// cross-building issue in case usage is 32bit
		ji.CPUTime = int64(usage.Utime.Sec) + int64(usage.Stime.Sec)
		// https://man7.org/linux/man-pages/man2/getrusage.2.html
		ji.ExtensionList[extension.JobInfoDefaultJSessionMaxRSS] = fmt.Sprintf("%d", usage.Maxrss)
		ji.ExtensionList[extension.JobInfoDefaultJSessionSwap] = fmt.Sprintf("%d", usage.Nswap)
		ji.ExtensionList[extension.JobInfoDefaultJSessionInBlock] = fmt.Sprintf("%d", usage.Inblock)
		ji.ExtensionList[extension.JobInfoDefaultJSessionOutBlock] = fmt.Sprintf("%d", usage.Oublock)
	}

	ji.ExtensionList[extension.JobInfoDefaultJSessionSystemTime] = fmt.Sprintf("%d", state.SystemTime().Milliseconds())
	ji.ExtensionList[extension.JobInfoDefaultJSessionUserTime] = fmt.Sprintf("%d", state.UserTime().Milliseconds())

	if state != nil && state.Success() {
		ji.State = drmaa2interface.Done
	} else {
		ji.State = drmaa2interface.Failed
	}

	if ji.ExitStatus != 0 {
		ji.State = drmaa2interface.Failed
	}

	ji.WallclockTime = time.Since(startTime)
	ji.ID = jobid
	ji.QueueName = ""

	return ji
}

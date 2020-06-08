package libdrmaa

import (
	"strconv"
	"strings"
	"time"

	"github.com/dgruber/drmaa"
	"github.com/dgruber/drmaa2interface"
)

// ConvertDRMAAJobInfoToDRMAA2JobInfo takes a drmaa v1 JobInfo object
// and converts it into a DRMAA2 JobInfo object.
func ConvertDRMAAJobInfoToDRMAA2JobInfo(ji *drmaa.JobInfo) (info drmaa2interface.JobInfo) {
	if ji == nil {
		return info
	}
	info.ExitStatus = int(ji.ExitStatus())
	if ji.HasExited() {
		if info.ExitStatus == 0 {
			info.State = drmaa2interface.Done
		} else {
			info.State = drmaa2interface.Failed
		}
	} else {
		// assumes that this function is called when the job is in an end state
		info.State = drmaa2interface.Failed
	}
	info.ID = ji.JobID()
	info.TerminatingSignal = ji.TerminationSignal()
	if ji.HasSignaled() || ji.HasAborted() || info.TerminatingSignal != "" {
		info.State = drmaa2interface.Failed
	}
	// TODO map resource usage
	// map[acct_cpu:0.4506 acct_io:0.0000 acct_iow:0.0000 acct_maxvmem:0.0000 acct_mem:0.0000
	// cpu:0.4506 end_time:1590752894.0000 exit_status:1.0000 io:0.0000 iow:0.0000 maxvmem:0.0000
	// mem:0.0000 priority:0.0000 ru_idrss:0.0000 ru_inblock:0.0000 ru_ismrss:0.0000 ru_isrss:0.0000
	// ru_ixrss:0.0000 ru_majflt:0.0000 ru_maxrss:3600.0000 ru_minflt:208.0000 ru_msgrcv:0.0000
	// ru_msgsnd:0.0000 ru_nivcsw:29.0000 ru_nsignals:0.0000 ru_nswap:0.0000 ru_nvcsw:1.0000
	// ru_oublock:8.0000 ru_stime:0.2798 ru_utime:0.1709 ru_wallclock:0.0000 signal:0.0000
	// start_time:1590752894.0000 submission_time:1590752891.0000 vmem:0.0000]
	submissionTime, exists := ji.ResourceUsage()["submission_time"]
	if exists {
		info.SubmissionTime = ConvertUnixToTime(submissionTime)
	}
	startTime, exists := ji.ResourceUsage()["start_time"]
	if exists {
		info.DispatchTime = ConvertUnixToTime(startTime)
	}
	finishTime, exists := ji.ResourceUsage()["end_time"]
	if exists {
		info.FinishTime = ConvertUnixToTime(finishTime)
	}
	return info
}

// ConvertUnixToTime converts something like 1590752891.0000 to time.
// Some systems report ms since epoch others just seconds.
func ConvertUnixToTime(t string) time.Time {
	sinceEpoch, _ := strconv.ParseInt(strings.Split(t, ".")[0], 10, 64)
	// assuming before year 2049
	if sinceEpoch < 2500000000 {
		// should be seconds
		return time.Unix(sinceEpoch, 0)
	}
	// assume milliseconds
	seconds := int64(sinceEpoch / 1000) // convert to seconds
	msWithoutSeconds := int64(sinceEpoch - (seconds * 1000))
	return time.Unix(seconds, msWithoutSeconds*1000000)
}

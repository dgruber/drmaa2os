package extension

// JobInfo extensions for process backend (default session)

// Monitoring Session jobs

const (
	JobInfoDefaultMSessionProcessName    string = "name"
	JobInfoDefaultMSessionCPUAffinity    string = "cpu_affinity"
	JobInfoDefaultMSessionCPUUsage       string = "cpu_usage"
	JobInfoDefaultMSessionMemoryUsage    string = "memory_usage"
	JobInfoDefaultMSessionMemoryUsageRSS string = "memory_usage_rss"
	JobInfoDefaultMSessionMemoryUsageVMS string = "memory_usage_vms"
	JobInfoDefaultMSessionCommandLine    string = "commandline"
	JobInfoDefaultMSessionWorkingDir     string = "workingdir"
)

// Job Session jobs

const (
	JobInfoDefaultJSessionMaxRSS     string = "ru_maxrss"
	JobInfoDefaultJSessionSwap       string = "ru_swap"
	JobInfoDefaultJSessionInBlock    string = "ru_inblock"
	JobInfoDefaultJSessionOutBlock   string = "ru_outblock"
	JobInfoDefaultJSessionSystemTime string = "system_time_ms"
	JobInfoDefaultJSessionUserTime   string = "user_time_ms"
)

// JobInfo extensions for Kubernetes backend (kubernetes session)

const (
	// JobInfoK8sJSessionJobOutput refers to the output of the job
	JobInfoK8sJSessionJobOutput string = "output"
)

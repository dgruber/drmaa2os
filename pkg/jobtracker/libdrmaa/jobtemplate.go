package libdrmaa

import (
	"fmt"
	"strings"

	"github.com/dgruber/drmaa"
	"github.com/dgruber/drmaa2interface"
)

// ConvertDRMAAJobTemplateToDRMAA2JobTemplate transforms a C DRMAA job template into
// a Go DRMAA2 job template.
func ConvertDRMAAJobTemplateToDRMAA2JobTemplate(jt *drmaa.JobTemplate) (drmaa2interface.JobTemplate, error) {
	if jt == nil {
		return drmaa2interface.JobTemplate{}, fmt.Errorf("job template is nil")
	}
	var t drmaa2interface.JobTemplate

	t.RemoteCommand, _ = jt.RemoteCommand()
	t.Args, _ = jt.Args()
	t.InputPath, _ = jt.InputPath()
	t.OutputPath, _ = jt.OutputPath()
	t.ErrorPath, _ = jt.ErrorPath()
	t.JoinFiles, _ = jt.JoinFiles()
	t.Email, _ = jt.Email()
	t.JobName, _ = jt.JobName()
	t.NativeSpecification, _ = jt.NativeSpecification()

	if submissionState, err := jt.JobSubmissionState(); err == nil && submissionState == drmaa.HoldState {
		t.SubmitAsHold = true
	}

	// job environment
	if block, err := jt.BlockEmail(); err == nil && block == true {
		t.Email = nil
		t.EmailOnStarted = false
		t.EmailOnTerminated = false
	}

	/*
	   if deadline, err := jt.DeadlineTime(); err == nil && deadline != time.Duration{} {
	      time.DeadlineTime = time.Now().Add(deadline)
	   }
	*/

	if env, err := jt.Env(); err == nil && len(env) > 0 {
		t.JobEnvironment = make(map[string]string, len(env))
		for _, v := range env {
			splittedEnv := strings.Split(v, "=")
			if len(splittedEnv) == 2 {
				t.JobEnvironment[splittedEnv[0]] = splittedEnv[1]
			}
		}
	}

	/*
	   func (jt *JobTemplate) HardRunDurationLimit() (deadlineTime time.Duration, err error)
	   func (jt *JobTemplate) HardWallclockTimeLimit() (deadlineTime time.Duration, err error)
	   func (jt *JobTemplate) NativeSpecification() (string, error)
	*/

	return t, nil
}

// ConvertDRMAA2JobTemplateToDRMAAJobTemplate transforms a Go DRMAA2 job template into
// a C drmaa job template.
func ConvertDRMAA2JobTemplateToDRMAAJobTemplate(jt drmaa2interface.JobTemplate, t *drmaa.JobTemplate) error {
	if jt.RemoteCommand != "" {
		t.SetRemoteCommand(jt.RemoteCommand)
	}
	if jt.Args != nil {
		t.SetArgs(jt.Args)
	}
	if jt.InputPath != "" {
		t.SetInputPath(":" + jt.InputPath)
	}
	if jt.OutputPath != "" {
		t.SetOutputPath(":" + jt.OutputPath)
	}
	if jt.ErrorPath != "" {
		t.SetErrorPath(":" + jt.ErrorPath)
	}
	t.SetJoinFiles(jt.JoinFiles)
	if jt.Email != nil {
		t.SetEmail(jt.Email)
	}
	if jt.JobName != "" {
		t.SetJobName(jt.JobName)
	}
	if jt.NativeSpecification != "" {
		t.SetNativeSpecification(jt.NativeSpecification)
	}

	if len(jt.JobEnvironment) > 0 {
		envs := make([]string, 0, len(jt.JobEnvironment))
		for k, v := range jt.JobEnvironment {
			envs = append(envs, fmt.Sprintf("%s=%s", k, v))
		}
		t.SetEnv(envs)
	}
	// missing:
	// BlockEmail / DeadlineTime / HardRunDurationLimit / HardWallclockTimeLimit
	return nil
}

package simpletracker

import (
	"errors"
	"fmt"
	"github.com/dgruber/drmaa2interface"
	"strconv"
	"strings"
)

// JobStore is an internal storage for jobs and job templates
// processed by the job tracker. Jobs are stored until Reap().
type JobStore struct {
	// jobids contains all known jobs in the system until they are reaped (Reap())
	// these are jobs, not array jobs and can be in format "1.1" or "1"
	jobids []string
	// running jobs
	// string is jobid and isArrayJob determines the type
	templates  map[string]drmaa2interface.JobTemplate
	jobs       map[string][]InternalJob
	isArrayJob map[string]bool
}

// NewJobStore returns a new in memory job store for jobs.
func NewJobStore() *JobStore {
	return &JobStore{
		jobids:     make([]string, 0, 512),
		templates:  make(map[string]drmaa2interface.JobTemplate),
		jobs:       make(map[string][]InternalJob),
		isArrayJob: make(map[string]bool),
	}
}

// SaveJob stores a job, the job submission template, and the process PID of
// the job in an internal job store.
func (js *JobStore) SaveJob(jobid string, t drmaa2interface.JobTemplate, pid int) {
	js.templates[jobid] = t
	js.jobids = append(js.jobids, jobid)
	js.jobs[jobid] = []InternalJob{
		{State: drmaa2interface.Running, PID: pid},
	}
}

// HasJob returns true if the job is saved in the job store.
func (js *JobStore) HasJob(jobid string) bool {
	_, exists := js.templates[jobid]
	return exists
}

// RemoveJob deletes all occurances of a job within the job storage.
// The jobid can be the identifier of a job or an job array.
func (js *JobStore) RemoveJob(jobid string) {
	isAJ, exits := js.isArrayJob[jobid]
	if exits && isAJ {
		jobids := make([]string, 0, len(js.jobids))
		for i := range js.jobids {
			if !strings.HasPrefix(js.jobids[i], jobid+".") {
				jobids = append(jobids, js.jobids[i])
			}
		}
		js.jobids = jobids
	} else {
		for i := range js.jobids {
			if js.jobids[i] == jobid {
				copy(js.jobids[i:], js.jobids[i+1:])
				js.jobids[len(js.jobids)-1] = ""
				js.jobids = js.jobids[:len(js.jobids)-1]
				break
			}
		}
	}
	delete(js.templates, jobid)
	delete(js.jobs, jobid)
	delete(js.isArrayJob, jobid)
}

func (js *JobStore) SaveArrayJob(arrayjobid string, pids []int, t drmaa2interface.JobTemplate, begin int, end int, step int) {
	pid := 0
	js.templates[arrayjobid] = t
	js.isArrayJob[arrayjobid] = true
	js.jobs[arrayjobid] = make([]InternalJob, 0, (end-begin)/step)

	for i := begin; i <= end; i += step {
		jobid := fmt.Sprintf("%s.%d", arrayjobid, i)
		js.jobids = append(js.jobids, jobid)
		js.jobs[arrayjobid] = append(js.jobs[arrayjobid], InternalJob{TaskID: i, State: drmaa2interface.Running, PID: pids[pid]})
		pid++
	}
}

func (js *JobStore) GetPID(jobid string) (int, error) {
	jobelements := strings.Split(jobid, ".")
	job, exists := js.jobs[jobelements[0]]
	if !exists {
		return -1, errors.New("Job does not exist")
	}
	var (
		taskid int
		err    error
	)
	if len(jobelements) > 1 {
		// is array job
		taskid, err = strconv.Atoi(jobelements[1])
		if err != nil {
			return -1, errors.New("TaskID within job ID is not a number")
		}
	}
	if taskid == 0 || taskid == 1 {
		return job[0].PID, nil
	}
	for task := range job {
		if job[task].TaskID == taskid {
			return job[task].PID, nil
		}
	}

	return -1, errors.New("TaskID not found in job array")
}

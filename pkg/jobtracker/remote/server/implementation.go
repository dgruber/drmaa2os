package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	genserver "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server/generated"
)

type JobTrackerImpl struct {
	jobTracker jobtracker.JobTracker
}

func NewJobTrackerImpl(jobTracker jobtracker.JobTracker) (*JobTrackerImpl, error) {
	return &JobTrackerImpl{
		jobTracker: jobTracker,
	}, nil
}

func (jti *JobTrackerImpl) AddArrayJob(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed reading body from addarrayjob request: %v\n", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	var aaj genserver.AddArrayJobJSONBody
	err = json.Unmarshal(body, &aaj)
	if err != nil {
		log.Printf("failed unmarshalling body from addarrayjob request: %v\n", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	step := 1
	if aaj.Step != nil {
		step = int(*aaj.Step)
	}
	maxParallel := 0
	if aaj.MaxParallel != nil {
		maxParallel = int(*aaj.MaxParallel)
	}
	id, err := jti.jobTracker.AddArrayJob(ConvertJobTemplateToDRMAA2(aaj.JobTemplate),
		int(aaj.Begin), int(aaj.End), step, maxParallel)
	var o genserver.AddArrayJobOutput
	if err != nil {
		o.Error = genserver.Error(err.Error())
	}
	o.JobID = genserver.JobID(id)
	out, _ := json.Marshal(o)
	success(w, out)
}

func (jti *JobTrackerImpl) AddJob(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed reading body from addjob request: %v\n", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	var jsonBody genserver.AddJobJSONBody
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		log.Printf("failed unmashalling body from addjob request: %v\n", err)
		http.Error(w, "can't unmarshal body", http.StatusBadRequest)
		return
	}
	id, err := jti.jobTracker.AddJob(ConvertJobTemplateToDRMAA2(genserver.JobTemplate(jsonBody)))

	var addJobOutput genserver.AddJobOutput
	addJobOutput.JobID = genserver.JobID(id)
	if err != nil {
		addJobOutput.Error = genserver.Error(err.Error())
	} else {
		addJobOutput.Error = genserver.Error("")
	}

	out, err := json.Marshal(addJobOutput)
	if err != nil {
		log.Printf("failed marshalling body for addjob response: %v\n", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	success(w, out)
}

func (jti *JobTrackerImpl) DeleteJob(w http.ResponseWriter, r *http.Request, params genserver.DeleteJobParams) {
	err := jti.jobTracker.DeleteJob(params.JobID)
	var response genserver.Error
	if err != nil {
		response = genserver.Error(err.Error())
	}
	out, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed marshalling body for deletejob response: %v\n", err)
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	success(w, out)
}

func (jti *JobTrackerImpl) JobControl(w http.ResponseWriter, r *http.Request, params genserver.JobControlParams) {
	err := jti.jobTracker.JobControl(params.JobID, string(params.Action))
	var response genserver.Error
	if err != nil {
		response = genserver.Error(err.Error())
	}
	out, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed marshalling body for jobcontrol response: %v\n", err)
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	success(w, out)
}

func (jti *JobTrackerImpl) JobInfo(w http.ResponseWriter, r *http.Request, params genserver.JobInfoParams) {
	var output genserver.JobInfoOutput
	ji, err := jti.jobTracker.JobInfo(params.JobID)
	if err != nil {
		output.Error = genserver.Error(err.Error())
	}
	output.JobInfo = ConvertJobInfo(ji)
	out, err := json.Marshal(output)
	if err != nil {
		log.Printf("failed marshalling body for jobinfo response: %v\n", err)
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	success(w, out)
}

func (jti *JobTrackerImpl) JobState(w http.ResponseWriter, r *http.Request, params genserver.JobStateParams) {
	state, substate, err := jti.jobTracker.JobState(params.JobID)
	result := genserver.JobStateOutput{
		JobState:    ConvertJobState(state.String()),
		JobSubState: genserver.JobSubState(substate),
	}
	out, err := json.Marshal(result)
	if err != nil {
		log.Printf("failed marshalling body for jobstate response: %v\n", err)
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	success(w, out)
}

func (jti *JobTrackerImpl) ListArrayJobs(w http.ResponseWriter, r *http.Request, params genserver.ListArrayJobsParams) {
	jobs, err := jti.jobTracker.ListArrayJobs(params.ArrayJobID)
	jobids := make([]genserver.JobID, 0, len(jobs))
	for _, job := range jobs {
		jobids = append(jobids, genserver.JobID(job))
	}
	out, err := json.Marshal(jobids)
	if err != nil {
		log.Printf("failed marshalling body for listarrayjobs response: %v\n", err)
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	success(w, out)
}

func (jti *JobTrackerImpl) ListJobCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := jti.jobTracker.ListJobCategories()
	out, err := json.Marshal(cats)
	if err != nil {
		log.Printf("failed marshalling body for listjobcategories response: %v\n", err)
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	success(w, out)
}

func (jti *JobTrackerImpl) ListJobs(w http.ResponseWriter, r *http.Request, params genserver.ListJobsParams) {
	jobs, err := jti.jobTracker.ListJobs()
	jobids := make([]genserver.JobID, 0, len(jobs))
	for _, job := range jobs {
		jobids = append(jobids, genserver.JobID(job))
	}
	out, err := json.Marshal(jobids)
	if err != nil {
		log.Printf("failed marshalling body for listjobs response: %v\n", err)
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	success(w, out)
}

func success(w http.ResponseWriter, out []byte) {
	w.Header().Set("Content-Type", "json")
	w.WriteHeader(200)
	w.Write(out)
}

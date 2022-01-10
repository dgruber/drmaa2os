package client

import (
	"context"
	"fmt"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/helper"
	genclient "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client/generated"
)

type ClientJobTracker struct {
	client genclient.ClientWithResponsesInterface
}

// init registers the remote client tracker at the SessionManager
func init() {
	drmaa2os.RegisterJobTracker(drmaa2os.RemoteSession, NewAllocator())
}

// New creates a new remote client job tracker.
func New(jobSessionName string, params ClientTrackerParams) (*ClientJobTracker, error) {
	if params.Server == "" {
		params.Server = "localhost:32321"
	}
	opts := make([]genclient.ClientOption, 0)
	if params.Path != "" {
		opts = []genclient.ClientOption{genclient.WithBaseURL(params.Server + params.Path)}
	}
	for _, v := range params.Opts {
		opts = append(opts, v)
	}
	client, err := genclient.NewClientWithResponses(
		params.Server,
		opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote client: %v", err)
	}

	return &ClientJobTracker{
		client: client,
	}, nil
}

func (c *ClientJobTracker) ListJobs() ([]string, error) {
	resp, err := c.client.ListJobsWithResponse(context.Background(),
		&genclient.ListJobsParams{})
	if err != nil || resp == nil {
		return nil, fmt.Errorf("failed listing jobs from remote: %v", err)
	}
	if resp.JSON200 == nil {
		return []string{}, nil
	}
	out := make([]string, 0, len(*resp.JSON200))
	for _, v := range *resp.JSON200 {
		out = append(out, string(v))
	}
	return out, nil
}

func (c *ClientJobTracker) AddJob(template drmaa2interface.JobTemplate) (string, error) {
	body := genclient.AddJobJSONRequestBody(ConvertJobTemplate(template))
	resp, err := c.client.AddJobWithResponse(context.Background(), body)
	if err != nil || resp == nil {
		return "", fmt.Errorf("failed adding job to remote: %v", err)
	}
	if resp.JSON200 == nil {
		return "", fmt.Errorf("failed adding job to remote")
	}
	if resp.JSON200.Error != "" {
		return string(resp.JSON200.JobID), fmt.Errorf("add job execution failed: %s", resp.JSON200.Error)
	}
	return string(resp.JSON200.JobID), nil
}

func (c *ClientJobTracker) AddArrayJob(jt drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error) {
	var body genclient.AddArrayJobJSONRequestBody

	body.JobTemplate = ConvertJobTemplate(jt)
	body.Begin = int64(begin)
	body.End = int64(end)
	if step != 0 {
		s := int64(step)
		body.Step = &s
	}
	if maxParallel != 0 {
		p := int64(maxParallel)
		body.MaxParallel = &p
	}
	resp, err := c.client.AddArrayJobWithResponse(context.Background(), body)
	if err != nil || resp == nil {
		return "", fmt.Errorf("failed adding job array to remote: %v", err)
	}
	if resp.JSON200 == nil {
		return "", fmt.Errorf("failed adding array job to remote")
	}
	if resp.JSON200.Error != "" {
		return string(resp.JSON200.JobID), fmt.Errorf("add array job execution failed: %s", resp.JSON200.Error)
	}
	return string(resp.JSON200.JobID), nil
}

func (c *ClientJobTracker) ListArrayJobs(arrayjobid string) ([]string, error) {
	resp, err := c.client.ListArrayJobsWithResponse(context.Background(), &genclient.ListArrayJobsParams{
		ArrayJobID: arrayjobid,
	})
	if err != nil || resp == nil {
		return nil, fmt.Errorf("failed listing array jobs from remote: %v", err)
	}
	if resp.JSON200 == nil {
		return []string{}, nil
	}
	out := make([]string, 0, len(*resp.JSON200))
	for _, v := range *resp.JSON200 {
		out = append(out, string(v))
	}
	return out, nil
}

func (c *ClientJobTracker) JobState(jobid string) (drmaa2interface.JobState, string, error) {
	resp, err := c.client.JobStateWithResponse(context.Background(),
		&genclient.JobStateParams{JobID: jobid})
	if err != nil || resp == nil {
		return drmaa2interface.Undetermined, "", fmt.Errorf("failed requesting job state: %v", err)
	}
	if resp.JSON200 == nil {
		return drmaa2interface.Undetermined, "", fmt.Errorf("failed requesting job state from remote: %v", err)
	}
	return ConvertJobStateToDRMAA2(string(resp.JSON200.JobState)), string(resp.JSON200.JobSubState), nil
}

func (c *ClientJobTracker) JobInfo(jobid string) (drmaa2interface.JobInfo, error) {
	resp, err := c.client.JobInfoWithResponse(context.Background(),
		&genclient.JobInfoParams{JobID: jobid})
	if err != nil || resp == nil {
		return drmaa2interface.JobInfo{}, fmt.Errorf("failed requesting job info: %v", err)
	}
	if resp.JSON200 == nil {
		return drmaa2interface.JobInfo{}, fmt.Errorf("failed requesting job info from remote")
	}
	if resp.JSON200.Error != "" {
		err = fmt.Errorf("failed requesting job info from remote: %v", err)
	}
	return ConvertJobInfoToDRMAA2(resp.JSON200.JobInfo), err
}

func (c *ClientJobTracker) JobControl(jobid, action string) error {
	resp, err := c.client.JobControlWithResponse(context.Background(),
		&genclient.JobControlParams{
			JobID:  jobid,
			Action: genclient.JobControlParamsAction(action),
		})
	if err != nil || resp == nil {
		return fmt.Errorf("failed changing job state: %v", err)
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("failed changing job state")
	}
	if *resp.JSON200 != "" {
		return fmt.Errorf("failed changing job state: %s", *resp.JSON200)
	}
	return nil
}

// Wait until the job has a certain DRMAA2 state or return an error if the state
// is unreachable.
func (p *ClientJobTracker) Wait(jobid string, timeout time.Duration, states ...drmaa2interface.JobState) error {
	// this is not specified in the OpenAPI spec for JobTracker as
	// we can rely on the helper function querying the job state
	// for now.
	return helper.WaitForStateWithInterval(p, 200*time.Millisecond, jobid, timeout, states...)
}

// DeleteJob removes a finished job from remote.
func (c *ClientJobTracker) DeleteJob(jobid string) error {
	resp, err := c.client.DeleteJobWithResponse(context.Background(), &genclient.DeleteJobParams{
		JobID: jobid,
	})
	if err != nil || resp == nil {
		return fmt.Errorf("failed deleting job from remote: %v", err)
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("failed deleting job from remote")
	}
	if *resp.JSON200 != "" {
		return fmt.Errorf("failed deleting job from remote: %s", *resp.JSON200)
	}
	return nil
}

// ListJobCategories returns all job categories from remote.
func (c *ClientJobTracker) ListJobCategories() ([]string, error) {
	resp, err := c.client.ListJobCategoriesWithResponse(context.Background())
	if err != nil || resp == nil {
		return nil, fmt.Errorf("failed listing jobs from remote: %v", err)
	}
	// when there are no job categories return an empty string slice
	if resp.JSON200 == nil {
		return []string{}, nil
	}
	return *resp.JSON200, nil
}

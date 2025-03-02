package simpletracker

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dgruber/drmaa2interface"
	bolt "go.etcd.io/bbolt"
)

const JobIDsStorageKey string = "JobIDsStorageKey"
const JobTemplatesStorageKey string = "JobTemplatesStorageKey"
const JobStorageKey string = "JobStorageKey"
const IsArrayJobStorageKey string = "IsArrayJobStorageKey"
const HighestJobIDStorageKey string = "HighestJobIDStorageKey"
const JobInfoStorageKey string = "JobInfoStorageKey"

// PersistentJobStorage is an internal storage for jobs and job templates
// processed by the job tracker. Jobs are stored until Reap().
// Locking must be done externally.
type PersistentJobStorage struct {
	//jobsession string
	// path to the DB file
	path string
	db   *bolt.DB
}

// NewPersistentJobStore returns a new job store which uses a file based DB
// to be persistent over process restarts. The PersistentJobStore implements
// the JobStorer interface.
func NewPersistentJobStore(path string) (*PersistentJobStorage, error) {
	jobstore := &PersistentJobStorage{
		//jobsession: jobsession,
		path: path,
	}

	var err error

	// allocate parent directory if it does not exist
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to create parent directory for job storage: %v\n", err)
		}
	}

	jobstore.db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialized boltdb for job storage: %v\n", err)
	}

	// ensure all buckets do exist
	err = jobstore.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(JobIDsStorageKey))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(JobTemplatesStorageKey))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(JobStorageKey))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(IsArrayJobStorageKey))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(HighestJobIDStorageKey))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(JobInfoStorageKey))
		if err != nil {
			return err
		}
		return nil
	})

	return jobstore, err
}

// SaveJob stores a job, the job submission template, and the process PID of
// the job in an internal job store.
func (js *PersistentJobStorage) SaveJob(jobid string, t drmaa2interface.JobTemplate, pid int) {

	err := js.db.Update(func(tx *bolt.Tx) error {
		db, err := tx.CreateBucketIfNotExists([]byte(JobIDsStorageKey))
		if err != nil {
			return err
		}
		err = db.Put([]byte(jobid), []byte(jobid))
		if err != nil {
			return fmt.Errorf("failed to save job: %v", err)
		}

		db, err = tx.CreateBucketIfNotExists([]byte(JobTemplatesStorageKey))
		if err != nil {
			return err
		}
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err = enc.Encode(t)
		if err != nil {
			return fmt.Errorf("failed to encode job template: %v", err)
		}
		err = db.Put([]byte(jobid), buffer.Bytes())
		if err != nil {
			return fmt.Errorf("failed to save job: %v", err)
		}

		db, err = tx.CreateBucketIfNotExists([]byte(JobStorageKey))
		if err != nil {
			return err
		}

		var jobbuffer bytes.Buffer
		enc = gob.NewEncoder(&jobbuffer)
		enc.Encode([]InternalJob{{State: drmaa2interface.Running, PID: pid}})
		err = db.Put([]byte(jobid), jobbuffer.Bytes())
		if err != nil {
			return fmt.Errorf("failed to save job: %v", err)
		}

		db, err = tx.CreateBucketIfNotExists([]byte(IsArrayJobStorageKey))
		if err != nil {
			return err
		}
		err = db.Put([]byte(jobid), []byte(fmt.Sprintf("%t", false)))
		if err != nil {
			return fmt.Errorf("failed to save job: %v", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("internal error: %v\n", err)
	}
}

// HasJob returns true if the job is saved in the job store.
func (js *PersistentJobStorage) HasJob(jobid string) bool {
	err := js.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(JobTemplatesStorageKey))
		if b == nil {
			// not found
			return fmt.Errorf("bucket with job templates not found")
		}
		template := b.Get([]byte(jobid))
		if template == nil {
			// not found, check in job list - might be an array job task
			bjid := tx.Bucket([]byte(JobIDsStorageKey))
			if bjid == nil {
				return fmt.Errorf("bucket with job ids not found")
			}
			jid := bjid.Get([]byte(jobid))
			if jid != nil {
				// job found
				return nil
			}
			return fmt.Errorf("jobid not found")
		}
		// job found
		return nil
	})
	return err == nil
}

func (js *PersistentJobStorage) IsArrayJob(jobid string) bool {
	err := js.db.View(func(tx *bolt.Tx) error {
		db := tx.Bucket([]byte(IsArrayJobStorageKey))
		if db == nil {
			return fmt.Errorf("bucket with name %s not found", IsArrayJobStorageKey)
		}
		isArrayJobDBEntry := db.Get([]byte(jobid))
		if isArrayJobDBEntry == nil {
			return fmt.Errorf("job %s is no array job", jobid)
		}
		// is array job
		return nil
	})
	if err != nil {
		return false
	}
	return true
}

// RemoveJob deletes all occurrences of a job within the job storage.
// The jobid can be the identifier of a job or a job array. In case
// of a job array it removes all tasks which belong to the array job.
func (js *PersistentJobStorage) RemoveJob(jobid string) {
	err := js.db.Update(func(tx *bolt.Tx) error {
		// is array job?
		db := tx.Bucket([]byte(IsArrayJobStorageKey))
		if db == nil {
			return fmt.Errorf("bucket with name %s not found", IsArrayJobStorageKey)
		}
		isArrayJob := false
		isArrayJobDBEntry := db.Get([]byte(jobid))
		if isArrayJobDBEntry != nil {
			isArrayJob = true
		}

		// delete all relevant jobs from that DB
		jobidsdb, err := tx.CreateBucketIfNotExists([]byte(JobIDsStorageKey))
		if err != nil {
			return err
		}

		if isArrayJob {
			// delete all jobs with that prefix
			// TODO optimize
			for _, somejobid := range js.GetJobIDs() {
				if strings.HasPrefix(somejobid, jobid+".") {
					// delete jobid
					jobidsdb.Delete([]byte(somejobid))
				}
			}
		} else {
			// delete only job id
			err = jobidsdb.Delete([]byte(jobid))
			if err != nil {
				return fmt.Errorf("failed to delete job %s: %v", jobid, err)
			}
		}

		db, err = tx.CreateBucketIfNotExists([]byte(JobTemplatesStorageKey))
		if err != nil {
			return err
		}
		db.Delete([]byte(jobid))

		db, err = tx.CreateBucketIfNotExists([]byte(JobStorageKey))
		if err != nil {
			return err
		}
		db.Delete([]byte(jobid))

		db, err = tx.CreateBucketIfNotExists([]byte(IsArrayJobStorageKey))
		if err != nil {
			return err
		}
		db.Delete([]byte(jobid))

		return nil
	})

	if err != nil {
		log.Printf("unexpected internal error while deleting job %s: %v\n", jobid, err)
	}
}

func (js *PersistentJobStorage) saveJobTemplate(tx *bolt.Tx, jobid string, template drmaa2interface.JobTemplate) error {
	db, err := tx.CreateBucketIfNotExists([]byte(JobTemplatesStorageKey))
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(template)
	err = db.Put([]byte(jobid), buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to save job: %v", err)
	}
	return nil
}

func (js *PersistentJobStorage) saveIsArrayJobID(tx *bolt.Tx, jobid string, isArrayJob bool) error {
	db, err := tx.CreateBucketIfNotExists([]byte(IsArrayJobStorageKey))
	if err != nil {
		return err
	}
	err = db.Put([]byte(jobid), []byte(fmt.Sprintf("%t", isArrayJob)))
	if err != nil {
		return fmt.Errorf("failed to save array job flag: %v", err)
	}
	return nil
}

func (js *PersistentJobStorage) saveInternalJobs(tx *bolt.Tx, jobid string, internalJobs []InternalJob) error {
	db, err := tx.CreateBucketIfNotExists([]byte(JobStorageKey))
	if err != nil {
		return err
	}
	var jobbuffer bytes.Buffer
	enc := gob.NewEncoder(&jobbuffer)
	err = enc.Encode(internalJobs)
	if err != nil {
		return fmt.Errorf("failed to encode internal jobs: %v", err)
	}
	err = db.Put([]byte(jobid), jobbuffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to save job: %v", err)
	}
	return nil
}

func (js *PersistentJobStorage) getInternalJobs(tx *bolt.Tx, jobid string) ([]InternalJob, error) {
	db := tx.Bucket([]byte(JobStorageKey))
	if db == nil {
		return nil, fmt.Errorf("bucket %s does not exist", JobStorageKey)
	}
	jobs := db.Get([]byte(jobid))
	if jobs == nil {
		return nil, errors.New("Job does not exist")
	}
	buffer := bytes.NewBuffer(jobs)
	dec := gob.NewDecoder(buffer)
	var internalJobs []InternalJob
	err := dec.Decode(&internalJobs)
	return internalJobs, err
}

func (js *PersistentJobStorage) saveJobID(tx *bolt.Tx, jobid string) error {
	db, err := tx.CreateBucketIfNotExists([]byte(JobIDsStorageKey))
	if err != nil {
		return err
	}
	err = db.Put([]byte(jobid), []byte(jobid))
	if err != nil {
		return fmt.Errorf("failed to save job: %v", err)
	}
	return nil
}

// SaveArrayJob stores all process IDs of the tasks of an array job.
func (js *PersistentJobStorage) SaveArrayJob(arrayjobid string, pids []int,
	t drmaa2interface.JobTemplate, begin, end, step int) {
	pid := 0

	js.db.Update(func(tx *bolt.Tx) error {
		err := js.saveJobTemplate(tx, arrayjobid, t)
		if err != nil {
			return err
		}

		err = js.saveIsArrayJobID(tx, arrayjobid, true)
		if err != nil {
			return err
		}

		internalJobs := make([]InternalJob, 0)

		for i := begin; i <= end; i += step {
			jobid := fmt.Sprintf("%s.%d", arrayjobid, i)

			err = js.saveJobID(tx, jobid)
			if err != nil {
				return err
			}

			internalJobs = append(internalJobs,
				InternalJob{
					TaskID: i,
					State:  drmaa2interface.Queued,
					PID:    pids[pid],
				})
			pid++
		}
		err = js.saveInternalJobs(tx, arrayjobid, internalJobs)
		if err != nil {
			return err
		}
		return nil
	})

}

// SaveArrayJobPID stores the current PID of main process of the
// job array task.
func (js *PersistentJobStorage) SaveArrayJobPID(arrayjobid string, taskid, pid int) error {
	return js.db.Update(func(tx *bolt.Tx) error {
		internalJobs, err := js.getInternalJobs(tx, arrayjobid)
		if err != nil {
			return fmt.Errorf("could not get internal jobs for array job id %s: %v",
				arrayjobid, err)
		}
		for task := range internalJobs {
			if internalJobs[task].TaskID == taskid {
				internalJobs[task].PID = pid
				internalJobs[task].State = drmaa2interface.Running
				err = js.saveInternalJobs(tx, arrayjobid, internalJobs)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return errors.New("task not found")
	})
}

// GetPID returns the PID of a job or an array job task.
// It returns -1 and an error if the job is not known.
func (js *PersistentJobStorage) GetPID(jobid string) (int, error) {
	jobelements := strings.Split(jobid, ".")
	var jobidint int

	err := js.db.View(func(tx *bolt.Tx) error {
		job, err := js.getInternalJobs(tx, jobelements[0])
		if err != nil {
			return fmt.Errorf("Error getting job %s: %v",
				jobelements[0], err)
		}
		var (
			taskid int
		)
		if len(jobelements) > 1 {
			// is array job
			taskid, err = strconv.Atoi(jobelements[1])
			if err != nil {
				return errors.New("TaskID within job ID is not a number")
			}
		}
		if taskid == 0 || taskid == 1 {
			jobidint = job[0].PID
			return nil
		}
		for task := range job {
			if job[task].TaskID == taskid {
				jobidint = job[task].PID
				return nil
			}
		}
		return errors.New("TaskID not found in job array")
	})

	if err != nil {
		return -1, err
	}

	return jobidint, nil
}

// GetJobIDs returns the IDs of all jobs.
func (js *PersistentJobStorage) GetJobIDs() []string {
	var jobs sync.Map

	err := js.db.View(func(tx *bolt.Tx) error {
		db := tx.Bucket([]byte(JobIDsStorageKey))
		if db == nil {
			return fmt.Errorf("bucket %s does not exist", JobIDsStorageKey)
		}
		db.ForEach(func(k []byte, v []byte) error {
			jobs.Store(string(k), string(v))
			return nil
		})
		return nil
	})
	if err != nil {
		log.Printf("internal error during getting job ids: %v", err)
	}

	jobids := make([]string, 0)
	jobs.Range(func(k interface{}, v interface{}) bool {
		jobids = append(jobids, k.(string))
		return true
	})

	return jobids
}

// GetArrayJobTaskIDs returns the IDs of all tasks of a job array.
func (js *PersistentJobStorage) GetArrayJobTaskIDs(arrayjobID string) []string {
	jobids := make([]string, 0)
	err := js.db.View(func(tx *bolt.Tx) error {
		internalJobs, err := js.getInternalJobs(tx, arrayjobID)
		if err != nil {
			return err
		}
		for _, job := range internalJobs {
			jobids = append(jobids, fmt.Sprintf("%s.%d", arrayjobID, job.TaskID))
		}
		return nil
	})
	if err != nil {
		return nil
	}
	// sort array jobs to have the same order as in-memory store
	sort.Strings(jobids)
	return jobids
}

func (js *PersistentJobStorage) NewJobID() string {
	highestjobid := ""
	err := js.db.Update(func(tx *bolt.Tx) error {
		db, err := tx.CreateBucketIfNotExists([]byte(HighestJobIDStorageKey))
		if err != nil {
			return err
		}
		// store highest job id - used for processes
		id := db.Get([]byte("highestjobid"))
		if id == nil {
			err := db.Put([]byte("highestjobid"), []byte("1"))
			if err != nil {
				return fmt.Errorf("failed to save job id as highest job id: %v", err)
			}
			highestjobid = "1"
		} else {
			// assume it is numerical like 1.1 or 1 -> +1
			jobid := strings.Split(string(id), ".")
			id, err := strconv.ParseInt(jobid[0], 10, 64)
			if err != nil {
				return fmt.Errorf("jobid not numerical: %v", err)
			}
			id++
			highestjobid = fmt.Sprintf("%d", id)
			err = db.Put([]byte("highestjobid"), []byte(highestjobid))
			if err != nil {
				return fmt.Errorf("failed to save job id as highest job id: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("failed to store highest job id: %v\n", err)
	}
	return highestjobid
}

func (js *PersistentJobStorage) Close() error {
	return js.db.Close()
}

func (js *PersistentJobStorage) GetJobTemplate(jobid string) (drmaa2interface.JobTemplate, error) {

	var jobTemplate drmaa2interface.JobTemplate

	err := js.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(JobTemplatesStorageKey))
		if b == nil {
			return fmt.Errorf("bucket with job templates not found")
		}

		template := b.Get([]byte(jobid))
		if template == nil {
			return fmt.Errorf("template for job %s not found", jobid)
		}

		buffer := bytes.NewBuffer(template)
		dec := gob.NewDecoder(buffer)
		return dec.Decode(&jobTemplate)
	})

	return jobTemplate, err
}

func (js *PersistentJobStorage) GetJobInfo(jobid string) (drmaa2interface.JobInfo, error) {
	var jobInfo drmaa2interface.JobInfo
	err := js.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(JobInfoStorageKey))
		if b == nil {
			return fmt.Errorf("bucket with job info not found")
		}

		info := b.Get([]byte(jobid))
		if info == nil {
			return fmt.Errorf("jobinfo for job %s not found", jobid)
		}

		buffer := bytes.NewBuffer(info)
		dec := gob.NewDecoder(buffer)
		return dec.Decode(&jobInfo)
	})
	return jobInfo, err
}

func (js *PersistentJobStorage) SaveJobInfo(jobid string, jobinfo drmaa2interface.JobInfo) error {
	return js.db.Update(func(tx *bolt.Tx) error {
		db := tx.Bucket([]byte(JobInfoStorageKey))
		if db == nil {
			return fmt.Errorf("bucket with job info not found")
		}
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(jobinfo)
		err := db.Put([]byte(jobid), buffer.Bytes())
		if err != nil {
			return fmt.Errorf("failed to save job info: %v", err)
		}
		return nil
	})
}

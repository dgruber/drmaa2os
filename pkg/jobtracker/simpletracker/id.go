package simpletracker

import (
	"fmt"
	"math"
)

var LastJobID int64

func GetNextJobID() string {
	if LastJobID == math.MaxInt64 {
		LastJobID = 0
	}
	LastJobID = LastJobID + 1
	return fmt.Sprintf("%d", LastJobID)
}

func SetJobID(jobid int64) {
	LastJobID = jobid
}

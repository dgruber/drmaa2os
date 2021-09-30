package drmaa2os

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"sync/atomic"

	"code.cloudfoundry.org/lager"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/dgruber/drmaa2os/pkg/storage"
	"github.com/dgruber/drmaa2os/pkg/storage/boltstore"
)

// atomicTrackers is the list of registered JobTrackers
var (
	trackerMutex   sync.Mutex
	atomicTrackers atomic.Value
)

// newRegisteredJobTracker creates a new JobTracker by calling the
// JobTracker creator which must be previously registered. Registering
// is done by importing the JobTracker packager where then the init()
// method is called. That decouples the JobTracker implementation from
// the rest of the code and only compiles dependencies which are required.
func (sm *SessionManager) newRegisteredJobTracker(jobSessionName string, params interface{}) (jobtracker.JobTracker, error) {
	jtMap := atomicTrackers.Load().(map[SessionType]jobtracker.Allocator)
	if jtMap == nil {
		return nil, errors.New("no JobTracker registered")
	}
	if _, exists := jtMap[sm.sessionType]; !exists {
		return nil, fmt.Errorf("JobTracker type %v not registered", sm.sessionType)
	}
	return jtMap[sm.sessionType].New(jobSessionName, params)
}

func makeSessionManager(dbpath string, st SessionType) (*SessionManager, error) {
	s := boltstore.NewBoltStore(dbpath)
	if err := s.Init(); err != nil {
		return nil, err
	}
	l := lager.NewLogger("sessionmanager")
	l.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))
	return &SessionManager{store: s, log: l, sessionType: st}, nil
}

func (sm *SessionManager) logErr(message string) error {
	return errors.New(message)
}

func (sm *SessionManager) create(t storage.KeyType, name string, contact string) error {
	if exists := sm.store.Exists(t, name); exists {
		return sm.logErr("Session already exists")
	}
	//if contact == "" {
	//	contact = name
	//}
	if err := sm.store.Put(t, name, contact); err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) delete(t storage.KeyType, name string) error {
	if err := sm.store.Delete(t, name); err != nil {
		return sm.logErr("Error while deleting")
	}
	return nil
}

// TryToSetContactString sets the contact string in the job tracker
// create params if create params has a contact string field.
func TryToSetContactString(createParams interface{}, contact string) error {
	// createParams must be pointer
	ps := reflect.ValueOf(createParams)
	if ps.Kind() != reflect.Ptr {
		return fmt.Errorf("createParams is not pointer")
	}
	createStruct := ps.Elem()
	// it is an interface value...
	//if createStruct.Kind() != reflect.Struct {
	//	return fmt.Errorf("createParams must be pointer to struct")
	//}
	f := createStruct.FieldByName("ContactString")
	if f.IsValid() {
		if f.Kind() == reflect.String {
			if f.CanSet() {
				f.SetString(contact)
			} else {
				return fmt.Errorf("can't set ContactString of create params")
			}
		} else {
			return fmt.Errorf("createParams has no ContactString value of kind string")
		}
	} else {
		return fmt.Errorf("ContactString is not a valid field in createParams")
	}
	return nil
}

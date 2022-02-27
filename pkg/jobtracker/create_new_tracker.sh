#!/bin/bash

# Helper script which sets up a new directory and code stubs
# for implementing a new job tacker. A job tracker is a simplified
# DRMAA2 backend implementation which is used by the framework
# to provide a DRMAA2 compatible interface.

# Usage:
#
# Let's assume you want to implement a new implementation for the
# workload manager called "test".
#
# Calling:
# ./create_new_tracker.sh test
# 
# - Creates a directory "testtracker"
# - Creates a file "testtracker/testtracker.go" with stubs for
#   a new Go type TestTracker which implements the JobTracker interface.

JOB_TRACKER_BACKEND_NAME=$1
JOB_TRACKER_DIR_NAME="${JOB_TRACKER_BACKEND_NAME}tracker"
JOB_TRACKER_NAME=`echo ${JOB_TRACKER_BACKEND_NAME:0:1} | tr '[a-z]' '[A-Z]'`${JOB_TRACKER_BACKEND_NAME:1}Tracker

function InstallImplTool() {
    if [ ! command -v impl &> /dev/null ]; then
        echo "impl command is not installed"
        echo "Installing impl for creating go interface implementation stubs"
        go install github.com/josharian/impl@latest
    fi
}

function CreateDirectoryAndFiles() {
    mkdir -p ${JOB_TRACKER_DIR_NAME}
    local gofile=${JOB_TRACKER_DIR_NAME}/${JOB_TRACKER_DIR_NAME}.go
    echo "package ${JOB_TRACKER_DIR_NAME}" > ${gofile}
    echo "" >> ${gofile}
    echo "import (" >> ${gofile}
    echo "      \"time\"" >> ${gofile}
    echo "      \"github.com/dgruber/drmaa2interface\"" >> ${gofile}
    echo ")" >> ${gofile}
    echo "" >> ${gofile}
    echo "type  ${JOB_TRACKER_NAME} struct {}" >> ${gofile}
    echo "" >> ${gofile}
    impl ${JOB_TRACKER_NAME} github.com/dgruber/drmaa2os/pkg/jobtracker.JobTracker >> ${JOB_TRACKER_DIR_NAME}/${JOB_TRACKER_DIR_NAME}.go
}

if [ "X${JOB_TRACKER_BACKEND_NAME}" == "X" ]; then 
    echo "Requires backend name (such as \"sarus\") as argument"
    exit 128
fi

echo "Using \"${JOB_TRACKER_BACKEND_NAME}\" as job tracker"
echo "Creating directory ${JOB_TRACKER_DIR_NAME}"
echo "Creating new job tracker ${JOB_TRACKER_NAME} implementation"

InstallImplTool
CreateDirectoryAndFiles

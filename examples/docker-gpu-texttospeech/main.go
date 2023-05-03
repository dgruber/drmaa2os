package main

import (
	"fmt"
	"os"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	// need to register docker backend
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"
)

func main() {

	sm, err := drmaa2os.NewDockerSessionManager("testdb.db")
	if err != nil {
		panic(err)
	}

	js, err := sm.OpenJobSession("tts-session")
	if err != nil {
		// it might not exist yet, try to create it (you can also invert the order)
		js, err = sm.CreateJobSession("tts-session", "docker")
		if err != nil {
			panic(err)
		}
	}

	localDir, _ := os.Getwd()

	// check here for CUDA and nvidia docker requirements of the host:
	// https://tts.readthedocs.io/en/latest/docker_images.html
	jt := drmaa2interface.JobTemplate{
		// Do not set command as it comes from the container (tts)
		Args:        []string{"--text", "Hello.", "--out_path", "/root/tts-output/hello.wav", "--use_cuda", "true"},
		JobCategory: "ghcr.io/coqui-ai/tts:v0.11.1", // latest doe not work for me
		StageInFiles: map[string]string{
			// format: "local": "container"
			localDir + "/tts-output": "/root/tts-output",
		},
		OutputPath: "/dev/stdout",
		ErrorPath:  "/dev/stderr",
		Extension: drmaa2interface.Extension{
			ExtensionList: map[string]string{
				"gpus": "all",
			},
		},
		JobEnvironment: map[string]string{
			"NVIDIA_VISIBLE_DEVICES":     "all",
			"NVIDIA_DRIVER_CAPABILITIES": "compute,utility",
			"NVIDIA_REQUIRE_CUDA":        "cuda>=11.0",
		},
	}

	job1, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	job1.WaitTerminated(drmaa2interface.InfiniteTime)
	if _, err := job1.GetJobInfo(); err != nil {
		panic(err)
	}

	if job1.GetState().String() == drmaa2interface.Done.String() {
		fmt.Println("Job finished successfully")
	} else {
		fmt.Println("Job finished with error")
	}

	// removing container
	job1.Reap()

	js.Close()
	sm.DestroyJobSession("tts-session")
}

package containerdtracker

type ContainerdTrackerParams struct {
	// ContainerdAddr is the address of the containerd daemon.
	ContainerdAddr string `json:"containerdAddr,omitempty"`
}

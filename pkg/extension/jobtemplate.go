package extension

// JobTemplate constants and extensions for Kubernetes backend

const (
	// JobTemplateK8sScheduler allows to specify a non-default Kubernetes scheduler for
	// the job when set.
	JobTemplateK8sScheduler string = "scheduler"
	// JobTemplateK8sLabels adds additional labels to the job object.
	// The value must be specified in the form: key=value,key=value,...
	JobTemplateK8sLabels string = "labels"
	// JobTemplateK8sPriviledged when set to TRUE runs the container of
	// the job in priviledged mode.
	JobTemplateK8sPrivileged string = "privileged"
	// JobTemplateK8sEnvFromSecret adds additional environment variables
	// specified in a secret to the container. The secret name must
	// be specified in the value. If multiple secrets are specified,
	// the names of the secrets must be colon separated.
	JobTemplateK8sEnvFromSecret string = "env-from-secret"
	// JobTemplateK8sEnvFromConfigMap adds additional environment variables
	// specified in a ConfigMap to the container. The ConfigMap name must
	// be specified in the value. If multiple ConfigMap are specified,
	// the names of the ConfigMaps must be colon separated.
	JobTemplateK8sEnvFromConfigMap string = "env-from-configmap"
	// JobTemplateK8sBasicSideCar when set to TRUE adds a basic sidecar
	// container to the job container. It stores the output of the
	// in a ConfigMap when the job is finished. This is only required
	// when the ConfigMap is consumed by someone (like a successor job).
	JobTemplateK8sBasicSideCar string = "DRMAA2_JOB_OUTPUT_IN_JOBINFO"
)

const (
	// JobTemplateK8sStageInAsSecretB64Prefix is a value prefix
	// in the JobTemplate StageInFiles map prefixing base64 encoded data
	// which is finally mounted into the job container as file defined
	// by the map key. The data (content of the file) itself is stored
	// as a Secret in the Kubernetes cluster.
	// Example: StageInFiles["/path/to/file"] = JobTemplateK8sStageInAsSecretB64Prefix + "some-base64-encoded-data"
	JobTemplateK8sStageInAsSecretB64Prefix string = "secret-data:"
	// JobTemplateK8sStageInAsConfigMapB64Prefix is a value prefix
	// in the JobTemplate StageInFiles map prefixing base64 encoded data
	// which is finally mounted into the job container as file defined
	// by the map key. The data (content of the file) itself is stored
	// as a ConfigMap in the Kubernetes cluster.
	// Example: StageInFiles["/path/to/file"] = JobTemplateK8sStageInAsConfigMapB64Prefix + "some-base64-encoded-data"
	JobTemplateK8sStageInAsConfigMapB64Prefix string = "configmap-data:"
	// JobTemplateK8sStageInFromStorageClassNamePrefix mounts a PVC derived from
	// a storage class name defined by to the specified path in the key of the map.
	// Example: StageInFiles["/storage"] = JobTemplateK8sStageInFromStorageClassNamePrefix + "some-storage-class-name"
	JobTemplateK8sStageInFromStorageClassNamePrefix string = "storageclass:"
	// JobTemplateK8sStageInFromPVCNamePrefix mounts path from the underlying
	// host into the container so that it can be accessed by the job.
	// Example: StageInFiles["/container/directory"] = JobTemplateK8sStageInFromPVCNamePrefix + "/host/directory"
	JobTemplateK8sStageInFromHostPathPrefix string = "hostpath:"
	// JobTemplateK8sStageInFromHostPathPrefix mounts an existing ConfigMap into the container
	// under the by key specified path.
	// Example: StageInFiles["/container/file"] = JobTemplateK8sStageInFromHostPathPrefix + "name-of-configmap-with-data"
	JobTemplateK8sStageInFromConfigMapPrefix string = "configmap:"
	// JobTemplateK8sStageInFromSecretPrefix mounts an existing Secret into the container
	// under the by key specified path.
	// Example: StageInFiles["/container/file"] = JobTemplateK8sStageInFromSecretPrefix + "name-of-secret-with-data"
	JobTemplateK8sStageInFromSecretPrefix string = "secret:"
	// JobTemplateK8sStageInFromPVCPrefix mounts an existing PVC into the container
	// under the by key specified path.
	// Example: StageInFiles["/container/dir"] = JobTemplateK8sStageInFromPVCPrefix + "name-of-pvc-with-data"
	JobTemplateK8sStageInFromPVCPrefix string = "pvc:"
	// JobTemplateK8sStageInFromGCEDiskPrefix mounts an existing GCE disk into the container
	// under the by key specified path. The GCEPersistentDisk is mounted Read/Write assuming
	// a ext4 filesystem.
	// Example: StageInFiles["/container/dir"] = JobTemplateK8sStageInFromGCEDiskPrefix + "name-of-gce-disk-with-data"
	JobTemplateK8sStageInFromGCEDiskPrefix string = "gce-disk:"
	// JobTemplateK8sStageInFromGCEDiskReadOnlyPrefix mounts an existing GCE disk into the container
	// under the by key specified path. The GCEPersistentDisk is mounted Read-Only assuming
	// a ext4 filesystem.
	// Example: StageInFiles["/container/dir"] = JobTemplateK8sStageInFromGCEDiskReadOnlyPrefix + "name-of-gce-disk-with-data"
	JobTemplateK8sStageInFromGCEDiskReadOnlyPrefix string = "gce-disk-read:"
	// JobTemplateK8sStageInFromNFSVolume mount an existing NFS volume into the container
	// under the by key specified path. After JobTemplateK8sStageInFromNFSVolume the hostname (or IP)
	// and the path to the volume must be specified separated by a colon.
	// Example: StageInFiles["/container/dir"] = JobTemplateK8sStageInFromNFSVolume + "server:path/to/nfs/volume"
	JobTemplateK8sStageInFromNFSVolumePrefix string = "nfs:"
)

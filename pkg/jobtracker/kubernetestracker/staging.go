package kubernetestracker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/extension"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	k8sapi "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getJobStageInSecrets(jt drmaa2interface.JobTemplate) ([]*v1.Secret, error) {
	if jt.StageInFiles == nil {
		return nil, nil
	}
	secrets := make([]*v1.Secret, 0, 2)
	for k, v := range jt.StageInFiles {
		// secret-data creates a new secret with the base64 decoded content
		if strings.HasPrefix(v, extension.JobTemplateK8sStageInAsSecretB64Prefix) {
			content := strings.TrimPrefix(v, extension.JobTemplateK8sStageInAsSecretB64Prefix)
			decoded, err := base64.StdEncoding.DecodeString(content)
			if err != nil {
				return nil, fmt.Errorf("failed to base64 decode the secret: %v", err)
			}
			_, file := filepath.Split(k)
			secrets = append(secrets, &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: secretName(jt.JobName, k),
					Labels: map[string]string{
						"origin": "drmaa2os",
						"jobid":  jt.JobName,
					},
				},
				Data: map[string][]byte{
					file: decoded,
				},
			})
		}
	}
	return secrets, nil
}

func getJobStageInConfigMaps(jt drmaa2interface.JobTemplate) ([]*v1.ConfigMap, error) {
	if jt.StageInFiles == nil {
		return nil, nil
	}
	configmaps := make([]*v1.ConfigMap, 0, 2)
	for k, v := range jt.StageInFiles {
		if strings.HasPrefix(v, extension.JobTemplateK8sStageInAsConfigMapB64Prefix) {
			content := strings.TrimPrefix(v, extension.JobTemplateK8sStageInAsConfigMapB64Prefix)
			decoded, err := base64.StdEncoding.DecodeString(content)
			if err != nil {
				return nil, fmt.Errorf("failed to base64-decode the configmap: %v", err)
			}
			_, file := filepath.Split(k)
			if file == "" {
				return nil, fmt.Errorf("StageInFiles map has no file path set: %s", k)
			}
			configmaps = append(configmaps,
				&v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: configMapName(jt.JobName, k),
						Labels: map[string]string{
							"origin": "drmaa2os",
							"jobid":  jt.JobName,
						},
					},
					Data: map[string]string{
						file: string(decoded),
					},
				})
		}
	}
	return configmaps, nil
}

func getJobStageInPVCs(jt drmaa2interface.JobTemplate) ([]*v1.PersistentVolumeClaim, error) {
	pvcs := make([]*v1.PersistentVolumeClaim, 0, 2)
	for k, v := range jt.StageInFiles {
		if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromStorageClassNamePrefix) {
			storageClassName := strings.TrimPrefix(v, extension.JobTemplateK8sStageInFromStorageClassNamePrefix)
			quantity := resource.NewQuantity(100, resource.DecimalSI)

			pvcs = append(pvcs, &v1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: pvcName(jt.JobName, k),
				},
				Spec: v1.PersistentVolumeClaimSpec{
					StorageClassName: &storageClassName,
					AccessModes: []v1.PersistentVolumeAccessMode{
						v1.ReadWriteOnce,
					},
					Resources: v1.VolumeResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceStorage: *quantity,
						},
					},
				}})
		}
	}
	return pvcs, nil
}

// removeArtifacts deletes all created secrets and configmaps
func removeArtifacts(cs *kubernetes.Clientset, jt drmaa2interface.JobTemplate, namespace string) error {
	if jt.StageInFiles == nil {
		return nil
	}

	var err error
	secrets, secretCreateErr := getJobStageInSecrets(jt)
	if secretCreateErr != nil {
		err = secretCreateErr
	}
	for _, secret := range secrets {
		errDelete := cs.CoreV1().Secrets(namespace).Delete(context.Background(),
			secret.Name, k8sapi.DeleteOptions{})
		if err != nil {
			err = fmt.Errorf("%w %v", err, errDelete)
		}
	}

	configmaps, cmCreateErr := getJobStageInConfigMaps(jt)
	if cmCreateErr != nil {
		err = fmt.Errorf("%w %v", err, cmCreateErr)
	}
	for _, cm := range configmaps {
		errDelete := cs.CoreV1().ConfigMaps(namespace).Delete(context.Background(),
			cm.Name, k8sapi.DeleteOptions{})
		if err != nil {
			err = fmt.Errorf("%w %v", err, errDelete)
		}
	}
	return err
}

// removeArtifactsByJobID removes all objects stored along with the job object:
// - secrets (stagein)
// - configmaps (stagein)
// - job template configmap
func removeArtifactsByJobID(cs *kubernetes.Clientset, jobID, namespace string) error {
	// list secrets and delete those which match the label and jobID
	secretList, err := cs.CoreV1().Secrets(namespace).List(context.Background(),
		metav1.ListOptions{})
	for _, secret := range secretList.Items {
		if strings.HasPrefix(secret.Name, jobID+"-") {
			errDelete := cs.CoreV1().Secrets(namespace).Delete(context.Background(),
				secret.Name, k8sapi.DeleteOptions{})
			if err != nil {
				fmt.Println("secret delete error")
				err = fmt.Errorf("%w %v", err, errDelete)
			}
		}
	}
	// list configmaps and delete those which match the label and jobID
	configMapList, err := cs.CoreV1().ConfigMaps(namespace).List(context.Background(),
		metav1.ListOptions{})
	for _, cm := range configMapList.Items {
		if strings.HasPrefix(cm.Name, jobID+"-") {
			errDelete := cs.CoreV1().ConfigMaps(namespace).Delete(context.Background(),
				cm.Name, k8sapi.DeleteOptions{})
			if err != nil {
				fmt.Println("configmap delete error")
				err = fmt.Errorf("%w %v", err, errDelete)
			}
		}
	}
	// list pvcs and delete job related pvcs
	pvcList, err := cs.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(),
		metav1.ListOptions{})
	for _, pvc := range pvcList.Items {
		if strings.HasPrefix(pvc.Name, jobID+"-") {
			errDelete := cs.CoreV1().PersistentVolumeClaims(namespace).Delete(context.Background(),
				pvc.Name, k8sapi.DeleteOptions{})
			if err != nil {
				fmt.Println("pvc delete error")
				err = fmt.Errorf("%w %v", err, errDelete)
			}
		}
	}

	// explicitly remove job template configmap (deleted already above)
	cs.CoreV1().ConfigMaps(namespace).Delete(context.Background(),
		jobID+"-jobtemplate-configmap", k8sapi.DeleteOptions{})
	return err
}

func storeJobTemplateInConfigMap(cs *kubernetes.Clientset, jt drmaa2interface.JobTemplate, namespace string) error {
	// remove content of secrets
	for _, v := range jt.StageInFiles {
		if strings.HasPrefix(v, "secret-data") {
			jt.StageInFiles[v] = "<removed>"
		}
	}
	m, err := convertJobTemplateToStringMap(jt)
	if err != nil {
		return err
	}
	_, err = cs.CoreV1().ConfigMaps(namespace).Create(context.Background(),
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: jt.JobName + "-jobtemplate-configmap",
			},
			Data: m,
		}, k8sapi.CreateOptions{})
	return err
}

func getJobTemplateFromConfigMap(cs *kubernetes.Clientset, jobID, namespace string) (*drmaa2interface.JobTemplate, error) {
	cm, err := cs.CoreV1().ConfigMaps(namespace).Get(context.Background(), jobID+"-jobtemplate-configmap", k8sapi.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not find configmap %s: %w", jobID+"-jobtemplate-configmap", err)
	}
	jt, exists := cm.Data["JobTemplate"]
	if !exists {
		return nil, fmt.Errorf("JobTemplate does not exist in configmap %s", jobID+"-jobtemplate-configmap")
	}
	var jobTemplate drmaa2interface.JobTemplate
	err = json.Unmarshal([]byte(jt), &jobTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JobTemplate data from configmap %s: %w", jobID+"-jobtemplate-configmap", err)
	}
	return &jobTemplate, nil
}

func convertJobTemplateToStringMap(jt drmaa2interface.JobTemplate) (map[string]string, error) {
	j, err := json.Marshal(jt)
	if err != nil {
		return nil, err
	}
	return map[string]string{"JobTemplate": string(j)}, nil
}

func getConfigMapData(cs kubernetes.Interface, namespace, configMapName, field string) (string, error) {
	cm, err := cs.CoreV1().ConfigMaps(namespace).Get(context.Background(), configMapName, k8sapi.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("could not find configmap %s: %w", configMapName, err)
	}
	data, exists := cm.Data[field]
	if !exists {
		return "", fmt.Errorf("configmap %s does not contain data %s", configMapName, field)
	}
	return data, nil
}

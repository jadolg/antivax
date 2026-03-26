package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	admission "k8s.io/api/admission/v1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

const LinkerdInjectAnnotation = "linkerd.io/inject"

func MutateCronjobs(body []byte) ([]byte, error) {
	admReview := admission.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	ar := admReview.Request
	if ar == nil {
		return nil, fmt.Errorf("received admission review with nil request")
	}

	var cronjob batchv1.CronJob
	if err := json.Unmarshal(ar.Object.Raw, &cronjob); err != nil {
		return nil, fmt.Errorf("unable to unmarshal cronjob json object: %v", err)
	}

	annotations := cronjob.Spec.JobTemplate.Spec.Template.Annotations
	const annotationsPath = "/spec/jobTemplate/spec/template/metadata/annotations"
	return buildResponse(&admReview, ar, annotations, annotationsPath, "cronjob", cronjob.Name, cronjob.Namespace)
}

func MutateJobs(body []byte) ([]byte, error) {
	admReview := admission.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	ar := admReview.Request
	if ar == nil {
		return nil, fmt.Errorf("received admission review with nil request")
	}

	var job batchv1.Job
	if err := json.Unmarshal(ar.Object.Raw, &job); err != nil {
		return nil, fmt.Errorf("unable to unmarshal job json object: %v", err)
	}

	annotations := job.Spec.Template.Annotations
	const annotationsPath = "/spec/template/metadata/annotations"
	return buildResponse(&admReview, ar, annotations, annotationsPath, "job", job.Name, job.Namespace)
}

func buildResponse(admReview *admission.AdmissionReview, ar *admission.AdmissionRequest, annotations map[string]string, annotationsPath, kind, name, namespace string) ([]byte, error) {
	resp := admission.AdmissionResponse{
		Allowed: true,
		UID:     ar.UID,
	}
	pT := admission.PatchTypeJSONPatch
	resp.PatchType = &pT

	patches := []interface{}{}
	if _, present := annotations[LinkerdInjectAnnotation]; !present {
		log.WithField(kind, name).WithField("namespace", namespace).Infof("Patching %s", kind)
		if len(annotations) == 0 {
			// No annotations exist yet — add the entire map.
			patches = append(patches, map[string]interface{}{
				"op":    "add",
				"path":  annotationsPath,
				"value": map[string]string{LinkerdInjectAnnotation: "disabled"},
			})
		} else {
			// Other annotations exist — add only our key to avoid overwriting them.
			// "/" in a JSON Pointer key is encoded as "~1".
			patches = append(patches, map[string]interface{}{
				"op":    "add",
				"path":  annotationsPath + "/linkerd.io~1inject",
				"value": "disabled",
			})
		}
	}

	var err error
	resp.Patch, err = json.Marshal(patches)
	if err != nil {
		return nil, err
	}

	resp.Result = &metav1.Status{Status: "Success"}
	admReview.Response = &resp

	responseBody, err := json.Marshal(admReview)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

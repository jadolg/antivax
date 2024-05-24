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

	var err error
	var cronjob *batchv1.CronJob

	var responseBody []byte
	ar := admReview.Request
	resp := admission.AdmissionResponse{}

	if ar != nil {
		if err := json.Unmarshal(ar.Object.Raw, &cronjob); err != nil {
			return nil, fmt.Errorf("unable unmarshal cronjob json object %v", err)
		}

		resp.Allowed = true
		resp.UID = ar.UID
		pT := admission.PatchTypeJSONPatch
		resp.PatchType = &pT

		var op map[string]interface{}
		if _, present := cronjob.Spec.JobTemplate.Spec.Template.Annotations[LinkerdInjectAnnotation]; !present {
			log.WithField("cronjob", cronjob.Name).WithField("namespace", cronjob.Namespace).Info("Patch Cronjob")
			op = map[string]interface{}{
				"op":    "add",
				"path":  "/spec/jobTemplate/spec/template/metadata/annotations",
				"value": map[string]string{LinkerdInjectAnnotation: "disabled"},
			}
			resp.Patch, err = json.Marshal([]interface{}{op})
			if err != nil {
				return nil, err
			}
		} else {
			resp.Patch, err = json.Marshal([]interface{}{})
			if err != nil {
				return nil, err
			}
		}

		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp
		responseBody, err = json.Marshal(admReview)
		if err != nil {
			return nil, err
		}
	}

	return responseBody, nil
}

func MutateJobs(body []byte) ([]byte, error) {
	admReview := admission.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	var err error
	var job *batchv1.Job

	var responseBody []byte
	ar := admReview.Request
	resp := admission.AdmissionResponse{}

	if ar != nil {
		if err := json.Unmarshal(ar.Object.Raw, &job); err != nil {
			return nil, fmt.Errorf("unable unmarshal job json object %v", err)
		}
		resp.Allowed = true
		resp.UID = ar.UID
		pT := admission.PatchTypeJSONPatch
		resp.PatchType = &pT

		var op map[string]interface{}
		if _, present := job.Spec.Template.Annotations[LinkerdInjectAnnotation]; !present {
			log.WithField("job", job.Name).WithField("namespace", job.Namespace).Info("Patch Job")
			op = map[string]interface{}{
				"op":    "add",
				"path":  "/spec/template/metadata/annotations",
				"value": map[string]string{LinkerdInjectAnnotation: "disabled"},
			}
			resp.Patch, err = json.Marshal([]interface{}{op})
			if err != nil {
				return nil, err
			}
		} else {
			resp.Patch, err = json.Marshal([]interface{}{})
			if err != nil {
				return nil, err
			}
		}

		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp
		responseBody, err = json.Marshal(admReview)
		if err != nil {
			return nil, err
		}
	}

	return responseBody, nil
}

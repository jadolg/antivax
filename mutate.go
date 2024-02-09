package main

import (
	"fmt"
	admission "k8s.io/api/admission/v1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

func Mutate(body []byte) ([]byte, error) {
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

		op := map[string]interface{}{}
		if _, present := cronjob.Spec.JobTemplate.Annotations["linkerd.io/inject"]; !present {
			op = map[string]interface{}{
				"op":    "add",
				"path":  "/spec/jobTemplate/spec/template/metadata/annotations",
				"value": map[string]string{"linkerd.io/inject": "disabled"},
			}
		}
		resp.Patch, err = json.Marshal([]interface{}{op})

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

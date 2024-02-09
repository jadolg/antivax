package main

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
	"testing"
)

func TestMutateValidCronjob(t *testing.T) {
	b, err := os.ReadFile("testdata/injectable_cronjob.json")
	assert.NoError(t, err)

	response, err := Mutate(b)
	assert.NoError(t, err)

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"add","path":"/spec/jobTemplate/spec/template/metadata/annotations","value":{"linkerd.io/inject":"disabled"}}]`, string(rr.Patch))
}

package main

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
	"testing"
)

func TestMutate(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		kind     string
		expected string
	}{
		{
			name:     "When the cronjob needs to be injected then the patch should be the annotation",
			file:     "testdata/injectable_cronjob.json",
			kind:     "cronjob",
			expected: `[{"op":"add","path":"/spec/jobTemplate/spec/template/metadata/annotations","value":{"linkerd.io/inject":"disabled"}}]`,
		},
		{
			name:     "When the cronjob does not need to be injected then the patch should be empty",
			file:     "testdata/non_injectable_cronjob.json",
			kind:     "cronjob",
			expected: "[]",
		},
		{
			name:     "When the job needs to be injected then the patch should be the annotation",
			file:     "testdata/injectable_job.json",
			kind:     "job",
			expected: `[{"op":"add","path":"/spec/template/metadata/annotations","value":{"linkerd.io/inject":"disabled"}}]`,
		},
		{
			name:     "When the job does not need to be injected then the patch should be empty",
			file:     "testdata/non_injectable_job.json",
			kind:     "job",
			expected: "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := os.ReadFile(tt.file)
			assert.NoError(t, err)

			var response []byte
			if tt.kind == "cronjob" {
				response, err = MutateCronjobs(b)
			} else if tt.kind == "job" {
				response, err = MutateJobs(b)
			} else {
				t.Fatalf("unknown kind %s", tt.kind)
			}

			assert.NoError(t, err)

			r := v1beta1.AdmissionReview{}
			err = json.Unmarshal(response, &r)
			assert.NoError(t, err, "failed to unmarshal with error %s", err)

			rr := r.Response
			assert.Equal(t, tt.expected, string(rr.Patch))
		})
	}
}

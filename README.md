# Antivax (pun intended)

A MutatingWebhook that will avoid Linkerd injection on your CronJobs and Jobs.

## Dependencies

- [cert-manager](https://cert-manager.io/) 

## How to use

1. Apply the `deployment.yaml` file to your cluster. This will create the
`antivax` namespace and deploy the mutating webhook.

2. Add the `antivax: enabled` label to the namespaces you want to avoid
Linkerd injection.

## What if I want to allow injection in only one Cronjob/Job?

Then you'll need to manually add the annotation to that resource:

1. CronJob example:

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: testcronjob
  namespace: antivax
spec:
  schedule: "* * * * *" #	Run every minute
  jobTemplate:
    spec:
      template:
        metadata:
          annotations:
            linkerd.io/inject: enabled
        spec:
          containers:
            - name: testcronjob
              image: busybox:latest
              imagePullPolicy: IfNotPresent
              command:
                - /bin/sh
                - -c
                - date; echo Hello!
          restartPolicy: OnFailure
```

2. Job example:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: testjob
  namespace: antivax
spec:
  template:
    metadata:
      annotations:
        linkerd.io/inject: enabled
    spec:
      automountServiceAccountToken: false
      containers:
        - name: testjob
          image: busybox:latest
          imagePullPolicy: IfNotPresent
          command:
            - /bin/sh
            - -c
            - date; echo Hello!
      restartPolicy: OnFailure
```

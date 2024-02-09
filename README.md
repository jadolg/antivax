# Antivax

A MutatingWebhook that will avoid Linkerd injection on your CronJobs

## Dependencies

- [cert-manager](https://cert-manager.io/) 

## How to use

Apply the `deployment.yaml` file to your cluster. This will create the
`antivax` namespace and deploy the mutating webhook.

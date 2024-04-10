# Antivax (pun intended)

A MutatingWebhook that will avoid Linkerd injection on your CronJobs

## Dependencies

- [cert-manager](https://cert-manager.io/) 

## How to use

1. Apply the `deployment.yaml` file to your cluster. This will create the
`antivax` namespace and deploy the mutating webhook.

2. Add the `antivax: enabled` label to the namespaces you want to avoid
Linkerd injection.

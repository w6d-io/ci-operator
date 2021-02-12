# CI Operator

CI Operator is a kubernetes extension base on Tekton Pipeline that installs and runs on your Kubernetes cluster.<br>
By creates a `Play` resource it creates all tekton resources defined. <br>
CI Operator is an open-source project that aims to simplified the CI/CD pipeline creation by using Tekton pipeline

## Getting started

To get started, use the helm chart 
Only been tested on Helm 3

```
helm repo add w6dio https://charts.w6d.io
helm repo update
helm install ciop w6dio/ci-operator
```

## Understand CI Operator

- [Creating a Step](steps.md)

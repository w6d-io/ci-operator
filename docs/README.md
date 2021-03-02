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

- [Creating a LimitCI](limitci.md) resource
- [Creating a Step](steps.md) resource
- [Creating a Play](play.md) resource

## Configuration

```yaml
---
domain: "example.ci"
ingress:
  class: nginx
  issuer: letsencrypt-prod
workspaces:
  - name: values
    description: "Values file place holder"
    mountPath: /helm/values
  - name: config
    description: "Helm config folder"
    mountPath: /root/.config/helm
  - name: artifacts
    description: "Values artifacts place holder"
    mountPath: /artifacts
  - name: source
    description: "Values source place holder"
    mountPath: /source
Volume:
  name: ws
  volumeClaimTemplate:
    spec:
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 2Gi
cluster_role: bot-cx-role
hash:
  salt: wildcard
  min_length: 16
minio:
  host: mino.svc:9000
  access_key: ACCESSKEYSAMPLE
  secret_key: secretkeysample
  bucket: values
vault:
  host: vault.svc:8200

```

The `domain` and `ingress` part are mandatory

## TODO
- add vault feature to record all secret
- add toggle feature for minio whether is defined or not

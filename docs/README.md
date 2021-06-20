# CI Operator

CI Operator is a kubernetes CRD base on [Tekton Pipeline](https://github.com/tektoncd/pipeline).

CI Operator targets is to easily create pipeline CI/CD by using [Tekton Pipeline](https://github.com/tektoncd/pipeline).

By creates one `Play` resource CI Operator builds all tekton resources for the pipeline.

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
Volume: # volume will be added in the pipeline resource. This is what will be used for workspaces
  name: ws
  volumeClaimTemplate:
    spec:
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 2Gi
pipeline_prefix: p6e-cx
values:
  deploy: # template to use for deployment instead of the on include in the code
    name: ci-operator-template
    key: deploy-values.yaml
cluster_role: bot-cx-role # clusterrole to bind in the rolebinding for pipeline and deployment
hash: # is use to hash the project id use in the domain generation
  salt: wildcard
  min_length: 16
minio: # minio is use to put values.yaml generate by the ci operator and use for deployment
  host: mino.svc:9000
  access_key: ACCESSKEYSAMPLE
  secret_key: secretkeysample
  bucket: values
vault: # vault can be used for keep dockerconfigjson, git_token and / or kubeconfig those elements can be set in the play resource in the secret part
  host: vault.svc:8200
  token: token
```

The `domain` and `ingress` part are mandatory

## Environment variables

[Predefined environment variables](variables.md) are fixed for all step


# Steps

## Configuring a `Step`

### Generic `Step`

Add `ci.w6d.io/kind: generic` in the annotation to set a `Step` as a generic

#### Sample
```yaml
apiVersion: ci.w6d.io/v1alpha1
kind: Step
metadata:
  name: build-step-generic
  annotations:
    ci.w6d.io/kind: generic
    ci.w6d.io/task: build
    ci.w6d.io/order: "0"
step:
  name: build-and-push
  image: "gcr.io/kaniko-project/executor:latest"
  env:
    - name: DOCKER_CONFIG
      value: "/tekton/home/.docker"
  command:
    - /kaniko/executor
  args:
    - --single-snapshot
    - --snapshotMode=redo
    - --use-new-run
    - $(params.flags[*])
    - --dockerfile=$(resources.inputs.source.path)/$(params.DOCKERFILE)
    - --destination=$(params.IMAGE)
    - --context=$(resources.inputs.source.path)/$(params.CONTEXT)
```

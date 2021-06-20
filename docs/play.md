# Play

Play is the core of this controller. Once all [steps](steps.md) are created
You can start to create `play` resource that will build your CI/CD pipeline by configure and creating all tekton resources 

Currently all pipeline is running in the project id namespace prefixed by `p6e-cx-`

For instance if the project id is 42 your pipeline namespace will be `p6e-cx-42`

Before any action your pipeline namespace have to be created

The `play` resource must be created in this same pipeline namespace

## :warning: WARNING ##

The script field in task override the script part into the step with the same name as the task

example :

```yaml
---
apiVersion: ci.w6d.io/v1alpha1
kind: Step
metadata:
  name: unit-test-init-generic
  annotations:
    ci.w6d.io/kind: generic
    ci.w6d.io/task: unit-test
    ci.w6d.io/order: "0"
step:
  name: init
  image: "busybox:latest"
  script: |
    echo "init"
---
apiVersion: ci.w6d.io/v1alpha1
kind: Step
metadata:
  name: unit-test-generic
  annotations:
    ci.w6d.io/kind: generic
    ci.w6d.io/task: unit-test
    ci.w6d.io/order: "1"
step:
  name: unit-test
  image: "busybox:latest"
  script: |
    echo "test"
---
apiVersion: ci.w6d.io/v1alpha1
kind: Play
metadata:
  name: play-1-1
  namespace: p6e-cx-1
spec:
  name: nodejs-sample
  environment: develop
  project_id: 1
  pipeline_id: 1
  repo_url: https://github.com/w6d-io/nodejs-sample.git
  commit:
    sha: 3010508ce47519b9b7444dcd2d2961796c874cff
    ref: main
  tasks:
    - unit-test:
        script:
          - echo "script overridden"
```

here the script in the step resource named `unit-test-generic` will be overridden by the `play` resource part.

## Configuring a `Play`

for instance let try a minimal play resource is a generic step have been created

```yaml
apiVersion: ci.w6d.io/v1alpha1
kind: Play
metadata:
  name: play-1-1
  namespace: p6e-cx-1
spec:
  name: nodejs-sample
  environment: develop
  project_id: 1
  pipeline_id: 1
  repo_url: https://github.com/w6d-io/nodejs-sample.git
  commit:
    sha: 3010508ce47519b9b7444dcd2d2961796c874cff
    ref: main
  tasks:
    - build:
        docker:
          filepath: Dockerfile
          context: .
```

all field available

```yaml
apiVersion: ci.w6d.io/v1alpha1
kind: Play
metadata:
  name: play-1-1
  namespace: p6e-cx-1
spec:
  name: nodejs-sample
  environment: develop
  scope:
    name: default
  stack:
    language: js
    package: npm
  project_id: 1
  pipeline_id: 1
  repo_url: https://github.com/w6d-io/nodejs-sample.git
  docker_url: w6dio/nodejs-sample:3010508c-main
  commit:
    sha: 3010508ce47519b9b7444dcd2d2961796c874cff
    ref: main
  pipeline_namespace: p6e-cx-1
  domain: nodejs-sample.example.ci
  external: false
  expose: true
  tasks:
    - unit-tests:
        image: node:latest
        script:
          - echo "unit-test1"
          - echo "unit-test2"
          - echo "unit-test3"
    - build:
        docker:
          filepath: Dockerfile
          context: .
        variables:
          BUILDENV: build
    - deploy:
        variables:
          BUILDENV: build
  secret:
    git_token: abcdef1234567890abcdef1234567890abcdef
    .dockerconfigjson: '{"auths":{"https://index.docker.io/v1/":{"auth":"personaldockertoken"}}}'
```

## explanation



| field                    | Description                                                                                                                                                                               |
|--------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|`spec.name`               | it is used for service name in deployment and in the value of the label application                                                                                                       |
|`spec.environment`        | it is use to build the namespace string where the application will be deploy by deploy task. This field is mandatory if you have enable deploy task and the flag `spec.external` is false |
|`spec.pipeline_namespace` | this optional field is the namespace where all resource for the pipeline will be created. If empty the default one `<pipeline_prefix>-<project_id>` will be used.                         |


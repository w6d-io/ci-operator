# Play

Play is the core of this controller. Once all [steps](steps.md) created
You can start to create `play` resource that will build your CI/CD pipeline by configure and creating all tekton resources 

Currently all pipeline is running in the project id namespace prefixed by `p6e-cx-`

For instance if the project id is 42 your pipeline namespace will be `p6e-cx-42`

Before any action your pipeline namespace have to be created

The `play` resource must be created in this same pipeline namespace

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
  domain: nodejs-sample.example.ci
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
  dependencies:
    mongodb:
      variables:
        HOST: "$DB_HOST"
        PASSWORD: "$DB_PASSWORD"
        USERNAME: "$DB_USERNAME"
  secret:
    git_token: abcdef1234567890abcdef1234567890abcdef
    .dockerconfigjson: '{"auths":{"https://index.docker.io/v1/":{"auth":"personaldockertoken"}}}'
```

# Variables

Following variable is available into all step 

- `$(resources.inputs.source.path)` path where the repository is cloned
- `$(workspaces.artifacts.path)` path where to put or get artifacts
- `$(params.flags[*])` param of array type :warning: not available in script part

## Predefined environment

Alongside your environment variables defined into the `play` resource the following variable is available in each step

With `W6D` as prefix

### common from `play` resource

```shell
~ > env
W6D_PROJECT_NAME="ci-operator"
W6D_PROJECT_LANGUAGE="js"
W6D_PROJECT_PACKAGE="npm"
W6D_ENVIRONMENT="production"
W6D_PROJECT_ID="1"
W6D_PIPELINE_ID="1"
W6D_REPOSITORY_URL="https://github.com/w6d-io/nodejs-sample"
W6D_COMMIT_SHA="3010508ce47519b9b7444dcd2d2961796c874cff"
W6D_COMMIT_SHORT_SHA="3010508c"
W6D_COMMIT_REF_NAME="main"
W6D_EXPOSE_DOMAIN="app.example.io"
W6D_EXPOSE="true"
W6D_EXTERNAL="false"
W6D_DOCKER_URL="w6dio/nodejs-sample:latest"
```

### from config

```shell
W6D_CONFIG_DEFAULT_DOMAIN="wildcard.sh"
W6D_CONFIG_CLUSTER_ROLE="ci-op"
W6D_CONFIG_INGRESS_CLASS="nginx"
W6D_CONFIG_NAMESPACE="ci-system"
```

### from task

Variable name will contain the task name from `play` resource.

Example with tasks named `unit-test`

```shell
W6D_UNIT_TEST_IMAGE="node:latest"
```

Example with tasks named `build`

```shell
W6D_BUILD_CONTEXT="."
W6D_BUILD_DOCKERFILE="Dockerfile"
W6D_BUILD_IMAGE="docker.io/w6dio/nodejs-sample:latest"
```

Example with tasks named `deploy`

```shell
W6D_DEPLOY_NAMESPACE="production"
```

# CI OPERATOR

[![Check](https://github.com/w6d-io/ci-operator/workflows/Check/badge.svg)](https://github.com/w6d-io/ci-operator/releases/latest)
[![codecov](https://codecov.io/gh/w6d-io/ci-operator/branch/main/graph/badge.svg?token=OYXGUIEDAH)](https://codecov.io/gh/w6d-io/ci-operator)
[![GitHub report](https://goreportcard.com/badge/github.com/w6d-io/ci-operator)](https://goreportcard.com/report/github.com/w6d-io/ci-operator)
[![Current release](https://img.shields.io/github/release/w6d-io/ci-operator.svg)](https://github.com/w6d-io/ci-operator/releases/latest)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/w6d-io/ci-operator)
[![GitHub](https://img.shields.io/github/license/w6d-io/ci-operator?style=flat)](https://github.com/w6d-io/ci-operator/blob/v0.11.1/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/w6d-io/ci-operator.svg)](https://github.com/w6d-io/ci-operator/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/w6d-io/ci-operator.svg)](https://github.com/w6d-io/ci-operator/pulls)
 
CI Operator is a kubernetes CRD base on [Tekton Pipeline](https://github.com/tektoncd/pipeline).

CI Operator targets is to easily create pipeline CI/CD by using [Tekton Pipeline](https://github.com/tektoncd/pipeline).

By creates one `Play` resource CI Operator builds all tekton resources for the pipeline.
[docs](/docs)

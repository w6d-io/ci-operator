module github.com/w6d-io/ci-operator

go 1.15

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr => github.com/go-logr/zapr v0.4.0
	k8s.io/api => k8s.io/api v0.19.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.7
	k8s.io/client-go => k8s.io/client-go v0.19.7
	k8s.io/client-go/informers => k8s.io/client-go/informers v0.19.7
)

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/avast/retry-go v3.0.0+incompatible // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.1.1
	github.com/hashicorp/vault/api v1.0.4
	github.com/minio/minio-go/v6 v6.0.57
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/tektoncd/pipeline v0.18.1
	github.com/w6d-io/hook v0.1.1
	go.uber.org/zap v1.16.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/pkg v0.0.0-20200922164940-4bf40ad82aab
	sigs.k8s.io/controller-runtime v0.6.1
)

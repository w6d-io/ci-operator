module github.com/w6d-io/ci-operator

go 1.15

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.2.1
	github.com/go-logr/zapr => github.com/go-logr/zapr v0.3.0
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
	k8s.io/client-go/informers => k8s.io/client-go/informers v0.19.0
)

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.2.1
	github.com/hashicorp/vault/api v1.0.4
	github.com/minio/minio-go/v6 v6.0.57
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/tektoncd/pipeline v0.18.1
	go.uber.org/zap v1.15.0
	golang.org/x/sys v0.0.0-20201202213521-69691e467435 // indirect
	golang.org/x/tools v0.0.0-20201202200335-bef1c476418a // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/pkg v0.0.0-20200922164940-4bf40ad82aab
	sigs.k8s.io/controller-runtime v0.6.1
)

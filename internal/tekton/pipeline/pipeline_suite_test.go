/*
Copyright 2020 WILDCARD SA.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 02/03/2021
*/
package pipeline_test

import (
    tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
    ci "github.com/w6d-io/ci-operator/api/v1alpha1"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/rest"
    "path/filepath"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/envtest"
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestPipeline(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Pipeline Suite")
}

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

var _ = BeforeSuite(func(done Done) {
    By("bootstrapping test environment")
    testEnv = &envtest.Environment{
        ErrorIfCRDPathMissing: false,
        CRDDirectoryPaths: []string{
            filepath.Join("..", "..", "..", "config", "crd" ,"bases"),
            filepath.Join("..", "..", "..", "third_party", "tektoncd", "pipeline", "config"),
        },
    }

    var err error
    cfg, err = testEnv.Start()
    Expect(err).ToNot(HaveOccurred())
    Expect(cfg).ToNot(BeNil())

    err = ci.AddToScheme(scheme.Scheme)
    Expect(err).NotTo(HaveOccurred())

    err = tkn.AddToScheme(scheme.Scheme)
    Expect(err).NotTo(HaveOccurred())

    k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
    Expect(err).ToNot(HaveOccurred())
    Expect(k8sClient).ToNot(BeNil())

    close(done)
}, 60)

var _ = AfterSuite(func() {
    By("tearing down the test environment")
    err := testEnv.Stop()
    Expect(err).ToNot(HaveOccurred())
})

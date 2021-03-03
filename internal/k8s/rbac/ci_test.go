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
Created on 03/03/2021
*/
package rbac_test

import (
    "context"
    "github.com/w6d-io/ci-operator/internal"
    "github.com/w6d-io/ci-operator/internal/config"
    "github.com/w6d-io/ci-operator/internal/k8s/rbac"
    "github.com/w6d-io/ci-operator/internal/util"
    "k8s.io/client-go/kubernetes/scheme"
    ctrl "sigs.k8s.io/controller-runtime"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    ci "github.com/w6d-io/ci-operator/api/v1alpha1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("CI", func() {
    When("", func() {
        var (
            play *ci.Play
            c    rbac.CI
        )
        BeforeEach(func() {
            play = &ci.Play{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "test-create-1",
                    UID:       "uuid-uuid-uuid-uuid",
                },
                Spec: ci.PlaySpec{
                    ProjectID:   1,
                    PipelineID:  1,
                },
            }
            ns := &corev1.Namespace{
                ObjectMeta: metav1.ObjectMeta{
                    Name: util.GetCINamespacedName(rbac.Prefix, play).Namespace,
                },
            }
            _ = k8sClient.Create(context.TODO(), ns)
            err := config.New("testdata/config.yaml")
            Expect(err).To(Succeed())
            c = rbac.CI{
                WorkFlowStruct: internal.WorkFlowStruct{
                    Play: play,
                },
            }
        })
        AfterEach(func() {
        })
        Context("we try", func() {
            It("fails due to crossed namespace creation", func() {
                play.Namespace = "default"
                err := c.Create(context.TODO(), k8sClient, ctrl.Log)
                Expect(err).ToNot(Succeed())
                Expect(err.Error()).To(ContainSubstring("cross-namespace"))
            })
            It("succeed creation", func() {
                play.Namespace = "p6e-cx-1"
                c.Scheme = scheme.Scheme
                err := c.Create(context.TODO(), k8sClient, ctrl.Log)
                Expect(err).To(Succeed())
            })
        })
    })
})

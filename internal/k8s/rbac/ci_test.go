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
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/rbac"
	"github.com/w6d-io/ci-operator/internal/util"
	"k8s.io/apimachinery/pkg/runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
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
					Name: "test-create-1",
					UID:  "uuid-uuid-uuid-uuid",
				},
				Spec: ci.PlaySpec{
					ProjectID:  1,
					PipelineID: 1,
				},
			}
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: util.GetCINamespacedName(rbac.Prefix, play).Namespace,
				},
			}
			_ = k8sClient.Create(ctx, ns)
			err := config.New("testdata/config.yaml")
			Expect(err).To(Succeed())
			c = rbac.CI{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play:   play,
					Scheme: scheme,
				},
			}
		})
		AfterEach(func() {
		})
		Context("we try", func() {
			It("fails due to kind missing", func() {
				play.Namespace = "default"
				c.Scheme = runtime.NewScheme()
				err := c.Create(ctx, k8sClient, ctrl.Log)
				Expect(err).ToNot(Succeed())
				Expect(err.Error()).To(ContainSubstring("no kind is registered for the type v1alpha1.Play"))
			})
			It("succeed creation", func() {
				By("create namespace")
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: "p6e-cx",
					},
				}
				Expect(k8sClient.Create(ctx, ns)).To(Succeed())

				play.Namespace = ns.GetName()
				err := c.Create(ctx, k8sClient, ctrl.Log)
				Expect(err).To(Succeed())
			})
		})
	})
})

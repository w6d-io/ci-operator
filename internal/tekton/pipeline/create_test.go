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
	"context"

	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/tekton/pipeline"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	Context("check all methods", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("success parse", func() {
			p := pipeline.Pipeline{
				Play: &ci.Play{
					Spec: ci.PlaySpec{
						Tasks: []map[ci.TaskType]ci.Task{
							{
								ci.UnitTests: ci.Task{},
							},
							{
								ci.Build: ci.Task{
									Image:  "",
									Script: nil,
									Variables: map[string]string{
										"TEST": "test",
									},
									Docker: ci.Docker{},
								},
							},
							{
								ci.Deploy: ci.Task{},
							},
							{
								ci.IntegrationTests: ci.Task{},
							},
							{
								ci.Clean: ci.Task{},
							},
							{
								ci.E2ETests: ci.Task{},
							},
							{
								"test": ci.Task{},
							},
						},
					},
				},
			}
			err := p.Parse(ctrl.Log)
			Expect(err).To(Succeed())
			//Expect(err).To(Equal(""))
		})
		It("failed create", func() {
			p := pipeline.Pipeline{
				Play: &ci.Play{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-create-1",
						Namespace: "p6e-cx-1",
						UID:       "uuid-uuid-uuid-uuid",
					},
					Spec: ci.PlaySpec{
						Tasks: []map[ci.TaskType]ci.Task{
							{
								ci.UnitTests: ci.Task{},
							},
							{
								ci.Build: ci.Task{
									Image:  "",
									Script: nil,
									Variables: map[string]string{
										"TEST": "test",
									},
									Docker: ci.Docker{},
								},
							},
							{
								ci.Deploy: ci.Task{},
							},
							{
								ci.IntegrationTests: ci.Task{},
							},
							{
								ci.Clean: ci.Task{},
							},
							{
								ci.E2ETests: ci.Task{},
							},
						},
					},
				},
			}
			err := p.Create(context.TODO(), k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("cross-namespace"))

			By("set a nonexistent namespace")
			config.SetNamespace("p6e-cx-1")
			By("set project id")
			p.Play.Spec.ProjectID = 1
			By("set pipeline id")
			p.Play.Spec.PipelineID = 1
			p.Scheme = scheme

			err = p.Create(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring(`namespaces "p6e-cx-1" not found`))
		})
		It("succeed create", func() {
			err := config.New("testdata/config.yaml")
			Expect(err).To(Succeed())
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-1",
				},
			}
			err = k8sClient.Create(context.TODO(), ns)
			Expect(err).To(Succeed())
			p := pipeline.Pipeline{
				Scheme: scheme,
				Play: &ci.Play{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-create-1",
						Namespace: "p6e-cx-1",
						UID:       "uuid-uuid-uuid-uuid",
					},
					Spec: ci.PlaySpec{
						Name:       "test",
						ProjectID:  1,
						PipelineID: 1,
						Tasks: []map[ci.TaskType]ci.Task{
							{
								ci.Build: ci.Task{
									Image:  "",
									Script: nil,
									Variables: map[string]string{
										"TEST": "test",
									},
									Docker: ci.Docker{},
								},
							},
						},
					},
				},
			}
			err = p.Create(context.TODO(), k8sClient, ctrl.Log)
			Expect(err).To(Succeed())
			//Expect(err.Error()).To(ContainSubstring("cross-namespace"))
		})
	})
})

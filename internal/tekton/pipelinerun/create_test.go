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
package pipelinerun_test

import (
	"context"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/tekton/pipelinerun"
	"github.com/w6d-io/ci-operator/internal/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	Context("test all methods", func() {
		It("parses", func() {
			By("load config")
			Expect(config.New("testdata/config.yaml")).To(Succeed())

			p := pipelinerun.PipelineRun{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						Spec: ci.PlaySpec{
							Name: "test",
							Stack: ci.Stack{
								Language: "test",
							},
							ProjectID:  1,
							PipelineID: 1,
							Commit: ci.Commit{
								SHA: "test_test_test",
								Ref: "test",
							},
							Tasks: []map[ci.TaskType]ci.Task{
								{
									ci.UnitTests: ci.Task{
										Image: "test/test:test",
										Script: []string{
											"echo",
											"test",
										},
									},
									ci.Build: ci.Task{
										Variables: map[string]string{
											"TEST": "Test",
										},
										Image: "test/test:test",
										Script: []string{
											"echo",
											"test",
										},
									},
									ci.Deploy: ci.Task{
										Variables: map[string]string{
											"TEST": "Test",
										},
										Image: "test/test:test",
										Script: []string{
											"echo",
											"test",
										},
									},
									ci.IntegrationTests: ci.Task{
										Image: "test/test:test",
										Script: []string{
											"echo",
											"test",
										},
									},
									ci.Clean: ci.Task{},
									ci.E2ETests: ci.Task{
										Image: "test/test:test",
										Script: []string{
											"echo",
											"test",
										},
									},
								},
							},
							Vault: &ci.Vault{
								Secrets: map[ci.SecretKind]ci.VaultSecret{
									ci.KubeConfig: {
										Path: "/secrets",
									},
								},
							},
						},
					},
				},
			}
			err := p.Parse(ctrl.Log)
			Expect(err).To(Succeed())
		})
		It("failed creation", func() {
			p := pipelinerun.PipelineRun{
				WorkFlowStruct: internal.WorkFlowStruct{
					Scheme: scheme,
					Play: &ci.Play{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-create-1",
							Namespace: "default",
							UID:       "uuid-uuid-uuid-uuid",
						},
						Spec: ci.PlaySpec{
							Name:       "test",
							ProjectID:  1,
							PipelineID: 1,
							Stack: ci.Stack{
								Language: "test",
							},
							Commit: ci.Commit{
								SHA: "test_test_test",
								Ref: "test",
							},
							Tasks: []map[ci.TaskType]ci.Task{
								{
									ci.Build: ci.Task{
										Variables: map[string]string{
											"TEST": "Test",
										},
										Image: "test/test:test",
										Script: []string{
											"echo",
											"test",
										},
									},
								},
							},
						},
					},
				},
			}
			err := p.Create(context.TODO(), k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("cross-namespace owner references"))

			By("Create in nonexist namespace")
			p.Play.Namespace = "p6e-cx-30"
			p.Play.Spec.ProjectID = 30
			err = p.Create(context.TODO(), k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring(`namespaces "p6e-cx-30" not found`))
		})
		It("create", func() {
			By("create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-14",
				},
			}
			Expect(k8sClient.Create(context.TODO(), ns)).To(Succeed())

			By("create play resource")
			play := &ci.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-create-1",
					Namespace: "p6e-cx-14",
					UID:       "uuid-uuid-uuid-uuid",
				},
				Spec: ci.PlaySpec{
					Name:       "test",
					ProjectID:  14,
					PipelineID: 1,
					Stack: ci.Stack{
						Language: "test",
					},
					Commit: ci.Commit{
						SHA: "test_test_test",
						Ref: "test",
					},
					Tasks: []map[ci.TaskType]ci.Task{
						{
							ci.Build: ci.Task{
								Variables: map[string]string{
									"TEST": "Test",
								},
								Image: "test/test:test",
								Script: []string{
									"echo",
									"test",
								},
							},
						},
					},
				},
			}
			By("create pipelinerun")
			p := pipelinerun.PipelineRun{
				WorkFlowStruct: internal.WorkFlowStruct{
					Scheme: scheme,
					Play:   play,
				},
			}
			Expect(p.Create(context.TODO(), k8sClient, ctrl.Log)).To(Succeed())

			By("create namespace #15")
			play.Namespace = "p6e-cx-15"
			play.Spec.ProjectID = 15
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-15",
				},
			}
			Expect(k8sClient.Create(context.TODO(), ns)).To(Succeed())

			By("create fake pod")
			po := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-fake-15",
					Namespace: "p6e-cx-15",
					Labels:    util.GetCILabels(play),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "web",
							Image: "nginx",
							Ports: []corev1.ContainerPort{
								{
									Name:          "web",
									ContainerPort: 80,
									Protocol:      "TCP",
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, po)).To(Succeed())

			//Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "pod-fake-15", Namespace: "p6e-cx-15"}, po)).To(Succeed())
			//ctrl.Log.V(1).Info("DEBUG TEST", "POD content", fmt.Sprintf("%+v",
			//	util.GetObjectContain(po)))
			//ctrl.Log.V(1).Info("DEBUG TEST", "T PR content", fmt.Sprintf("%+v",
			//	p))

			By("throw delaying")
			Expect(p.Create(context.TODO(), k8sClient, ctrl.Log)).To(Succeed())
		})

		It("does", func() {
			p := pipelinerun.PipelineRun{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						Spec: ci.PlaySpec{
							Name:       "test",
							ProjectID:  1,
							PipelineID: 1,
							DockerURL:  "http://{}",
							Stack: ci.Stack{
								Language: "test",
							},
							Commit: ci.Commit{
								SHA: "test_test_test",
								Ref: "test",
							},
							Tasks: []map[ci.TaskType]ci.Task{
								{
									ci.Build: ci.Task{
										Variables: map[string]string{
											"TEST": "Test",
										},
										Image: "test/test:test",
										Script: []string{
											"echo",
											"test",
										},
									},
									"git-leaks": ci.Task{},
								},
							},
						},
					},
				},
			}
			err := p.Parse(ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("invalid character"))
		})
	})
})

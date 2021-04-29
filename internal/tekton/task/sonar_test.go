/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 29/04/2021
*/

package task_test

import (
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
	"github.com/w6d-io/ci-operator/internal/tekton/task"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	Context("Sonar", func() {
		It("execute Sonar", func() {
			var err error
			By("Create task")
			t := &task.Task{
				Index:  0,
				Client: k8sClient,
				Play: &ci.Play{
					Spec: ci.PlaySpec{
						Stack: ci.Stack{
							Language: "bash",
							Package:  "test",
						},
						Tasks: []map[ci.TaskType]ci.Task{
							{
								ci.Sonar: ci.Task{
									Script: ci.Script{
										"echo", "toto",
									},
								},
							},
						},
						Vault: &ci.Vault{
							Secrets: map[ci.SecretKind]ci.VaultSecret{
								secrets.KubeConfigKey: {},
							},
						},
					},
				},
			}

			By("load config")
			Expect(config.New("testdata/config.yaml")).To(Succeed())

			By("Set namespace")
			config.SetNamespace("p6e-cx-28")

			By("Create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-28",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			By("failed by no step")
			err = t.Sonar(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())

			By("Create step")
			step := &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "step-unit-test-1",
					Namespace: "p6e-cx-28",
					Annotations: map[string]string{
						ci.AnnotationPackage:  "test",
						ci.AnnotationLanguage: "bash",
						ci.AnnotationTask:     ci.Sonar.String(),
						ci.AnnotationOrder:    "0",
					},
				},
				Step: ci.StepSpec{
					Step: tkn.Step{
						Script: "echo test",
					},
				},
			}
			Expect(k8sClient.Create(ctx, step)).To(Succeed())

			Expect(t.Sonar(ctx, ctrl.Log)).To(Succeed())
		})
		It("Execute Create", func() {
			var err error
			By("build SonarTask")
			u := &task.SonarTask{
				Meta: task.Meta{
					Steps: []tkn.Step{
						{
							Script: "echo test",
						},
					},
					Play: &ci.Play{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "play-test-29-1",
							Namespace: "p6e-cx-29",
							UID:       "cf3e4129-6b00-410e-9d3d-774292a57bce",
						},
						Spec: ci.PlaySpec{
							Stack: ci.Stack{
								Language: "bash",
								Package:  "test",
							},
							Tasks: []map[ci.TaskType]ci.Task{
								{
									ci.Sonar: ci.Task{
										Script: ci.Script{
											"echo", "toto",
										},
									},
								},
							},
							Vault: &ci.Vault{
								Secrets: map[ci.SecretKind]ci.VaultSecret{
									secrets.KubeConfigKey: {},
								},
							},
						},
					},
					Scheme: scheme,
				},
			}

			By("Set config")
			Expect(config.New("testdata/config.yaml")).To(Succeed())

			By("Failed due to cross namespace")
			err = u.Create(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("cross-namespace owner references are disallowed"))

			By("set the right namespace")
			u.Play.Spec.PipelineID = 1
			u.Play.Spec.ProjectID = 29
			err = u.Create(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(Equal(`namespaces "p6e-cx-29" not found`))

			By("create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-29",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			Expect(u.Create(ctx, k8sClient, ctrl.Log)).To(Succeed())

		})
	})
})
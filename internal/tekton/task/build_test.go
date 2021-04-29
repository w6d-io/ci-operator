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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
	"github.com/w6d-io/ci-operator/internal/tekton/task"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("Task", func() {
	Context("Build", func() {
		It("execute build", func() {
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
								ci.Build: ci.Task{
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
			config.SetNamespace("p6e-cx-22")

			By("Create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-22",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			By("failed by no step")
			err = t.Build(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())

			By("Create step")
			step := &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "step-unit-test-1",
					Namespace: "p6e-cx-22",
					Annotations: map[string]string{
						ci.AnnotationPackage:  "test",
						ci.AnnotationLanguage: "bash",
						ci.AnnotationTask:     ci.Build.String(),
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

			Expect(t.Build(ctx, ctrl.Log)).To(Succeed())
		})
		It("Execute Create", func() {
			var err error
			By("build BuildTask")
			b := &task.BuildTask{
				Meta: task.Meta{
					Steps: []tkn.Step{
						{
							Script: "echo test",
						},
					},
					Play: &ci.Play{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "play-test-23-1",
							Namespace: "p6e-cx-23",
							UID:       "bb8fd459-bd40-43e3-b8c0-2fc63bc575a0",
						},
						Spec: ci.PlaySpec{
							Stack: ci.Stack{
								Language: "bash",
								Package:  "test",
							},
							Tasks: []map[ci.TaskType]ci.Task{
								{
									ci.Build: ci.Task{
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
				BuildDocker: true,
			}

			By("Set config")
			Expect(config.New("testdata/config.yaml")).To(Succeed())

			By("Failed due to cross namespace")
			err = b.Create(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("cross-namespace owner references are disallowed"))

			By("set the right namespace")
			b.Play.Spec.PipelineID = 1
			b.Play.Spec.ProjectID = 23
			err = b.Create(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(Equal(`namespaces "p6e-cx-23" not found`))

			By("create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-23",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			Expect(b.Create(ctx, k8sClient, ctrl.Log)).To(Succeed())

		})
	})
})

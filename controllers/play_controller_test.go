/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 20/04/2021
*/

package controllers_test

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/controllers"
	"github.com/w6d-io/ci-operator/internal/minio"
	"github.com/w6d-io/ci-operator/internal/util"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck/v1beta1"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Play controller", func() {
	Context("", func() {
		It("", func() {
			var err error
			namespace := "p6e-cx-99"

			By("create namespace #99")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())
			By("creating a new Play")
			play := &ci.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "play-99-1",
					Namespace: namespace,
					UID:       "2bbf1590-ef1d-4928-9673-d68f707c3827",
					Labels: map[string]string{
						"projectid":  "99",
						"pipelineid": "1",
					},
				},
				Spec: ci.PlaySpec{
					Name: "nodejs-sample",
					Stack: ci.Stack{
						Language: "js",
						Package:  "npm",
					},
					Environment: "prod",
					ProjectID:   99,
					PipelineID:  1,
					RepoURL:     "https://github.com/w6d-io/nodejs-sample",
					Commit: ci.Commit{
						SHA:       "3010508ce47519b9b7444dcd2d2961796c874cff",
						BeforeSHA: "0000000000000000000000000000000000000000",
						Ref:       "main",
						Message:   "init",
					},
					Domain:   "nodejs-sample.wildcard.sh",
					Expose:   true,
					External: false,
					Tasks: []map[ci.TaskType]ci.Task{
						{
							"unit-tests": ci.Task{
								Image: "busybox",
								Script: []string{
									"echo test",
								},
								Arguments: []string{
									"arg1",
									"arg2",
								},
								Variables: fields.Set{
									"VAR1": "value1",
									"VAR2": "value2",
								},
								Namespace: namespace,
								Docker: ci.Docker{
									Filepath: "Dockerfile",
									Context:  ".",
								},
								Annotations: map[string]string{
									ci.AnnotationKind: "unit-test",
								},
							},
							"git-leaks": ci.Task{},
						},
					},
					DockerURL: "docker.io/w6dio/nodejs-sample:latest",
				},
			}
			Expect(k8sClient.Create(ctx, play)).To(Succeed())

			By("create limit ci")
			limit := &ci.LimitCi{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "limit-test-1-1",
					Namespace: namespace,
				},
				Spec: ci.LimitCiSpec{
					Concurrent: 1,
				},
			}
			Expect(k8sClient.Create(ctx, limit)).To(Succeed())

			By("update play status")
			meta.SetStatusCondition(&play.Status.Conditions, metav1.Condition{
				Type:               "Creating",
				Status:             metav1.ConditionUnknown,
				LastTransitionTime: metav1.Time{Time: time.Now()},
				Reason:             "Testing",
				Message:            "unit test",
			})
			Expect(k8sClient.Status().Update(ctx, play)).To(Succeed())

			By("create types namespace")
			nn := types.NamespacedName{Namespace: namespace, Name: "play-99-1"}
			r := &controllers.PlayReconciler{
				Client: k8sClient,
				Log:    ctrl.Log,
				Scheme: scheme,
			}
			req := ctrl.Request{
				NamespacedName: nn,
			}
			_, err = r.Reconcile(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("minio connection failed"))

			By("Set minio")
			minio.Instance = &MockMinio{}

			_, err = r.Reconcile(ctx, req)
			Expect(err).To(HaveOccurred())
			//Expect(err.Error()).To(ContainSubstring("no step found for "))

			By("Play not found")
			req.Name = "play-99-2"
			_, err = r.Reconcile(ctx, req)
			Expect(err).ToNot(HaveOccurred())

			By("Create pipelinerun")
			pr := &tkn.PipelineRun{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pipeline-run-99-1",
					Namespace: namespace,
					Labels:    util.GetCILabels(play),
				},
				Spec: tkn.PipelineRunSpec{
					PipelineRef: &tkn.PipelineRef{
						Name: "test",
					},
				},
			}
			Expect(k8sClient.Create(ctx, pr)).To(Succeed())

			By("set pipelinerun status")
			pr.Status = tkn.PipelineRunStatus{
				Status: v1beta1.Status{
					Conditions: v1beta1.Conditions{
						{
							Type:               apis.ConditionReady,
							Status:             metav1.StatusFailure,
							LastTransitionTime: apis.VolatileTime{Inner: metav1.Time{Time: time.Now()}},
							Reason:             "no reason",
							Message:            "unit tests",
						},
					},
				},
			}
			Expect(k8sClient.Status().Update(ctx, pr))

			pr2 := &tkn.PipelineRun{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pipeline-run-99-2",
					Namespace: namespace,
					Labels:    util.GetCILabels(play),
				},
				Spec: tkn.PipelineRunSpec{
					PipelineRef: &tkn.PipelineRef{
						Name: "test",
					},
				},
			}
			Expect(k8sClient.Create(ctx, pr2)).To(Succeed())

			By("run play with pipelinerun exists")
			req.Name = "play-99-1"
			_, err = r.Reconcile(ctx, req)
			Expect(err).To(HaveOccurred())
			//Expect(err.Error()).To(Equal("xyz"))
		})
		It("", func() {
			r := &controllers.PlayReconciler{}
			Expect(r.GetStatus(ci.Succeeded)).To(Equal(metav1.ConditionTrue))
		})
	})
})

type MockMinio struct{}

func (m *MockMinio) PutFile(_ logr.Logger, _ string, _ string) error {
	return nil
}

func (m *MockMinio) PutString(_ logr.Logger, _ string, _ string) error {
	return nil
}

/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 31/12/2020
*/

package util_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis/duck/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

var _ = Describe("Util", func() {
	Context("", func() {
		It("Condition", func() {
			var c v1beta1.Conditions

			By("Empty conditions")
			Expect(util.Condition(c)).To(Equal(ci.State("---")))

			By("Condition false")
			c = v1beta1.Conditions{
				{
					Status: corev1.ConditionFalse,
				},
			}
			Expect(util.Condition(c)).To(Equal(ci.Failed))

			By("Condition true")
			c = v1beta1.Conditions{
				{
					Status: corev1.ConditionTrue,
				},
			}
			Expect(util.Condition(c)).To(Equal(ci.Succeeded))

			By("Condition unknown")
			c = v1beta1.Conditions{
				{
					Status: corev1.ConditionUnknown,
				},
			}
			Expect(util.Condition(c)).To(Equal(ci.Running))

			By("Condition cancel")
			c = v1beta1.Conditions{
				{
					Status: corev1.ConditionUnknown,
					Reason: "PipelineRunCancelled",
				},
			}
			Expect(util.Condition(c)).To(Equal(ci.Cancelled))

		})
		It("Message", func() {
			var c v1beta1.Conditions

			By("No Message")
			Expect(util.Message(c)).To(Equal(""))

			By("Set Message")
			c = v1beta1.Conditions{
				{
					Message: "Test",
				},
			}
			Expect(util.Message(c)).To(Equal("Test"))
		})
		It("IsPipelineRunning", func() {
			var pr tkn.PipelineRun

			pr.Status.Conditions = v1beta1.Conditions{
				{
					Status: corev1.ConditionFalse,
				},
			}
			Expect(util.IsPipelineRunning(pr)).To(Equal(false))

			pr.Status.Conditions = v1beta1.Conditions{
				{
					Status: corev1.ConditionUnknown,
				},
			}
			Expect(util.IsPipelineRunning(pr)).To(Equal(true))
		})
		It("InNamespace", func() {
			By("build play")
			play := &ci.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "play-tool-1",
					Namespace: "p6e-cx-1",
				},
				Spec: ci.PlaySpec{
					ProjectID: 1,
				},
			}
			Expect(util.InNamespace(play)).To(Equal(client.InNamespace("p6e-cx-1")))
		})
		It("GetCINamespacedName", func() {
			By("build play")
			play := &ci.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "play-tool-1",
					Namespace: "p6e-cx-1",
				},
				Spec: ci.PlaySpec{
					ProjectID:  1,
					PipelineID: 1,
				},
			}
			Expect(util.GetCINamespacedName("test", play).String()).To(Equal("p6e-cx-1/test-1-1"))
		})
		It("GetCINamespacedName2", func() {
			By("build play")
			play := &ci.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "play-tool-1",
					Namespace: "p6e-cx-1",
				},
				Spec: ci.PlaySpec{
					ProjectID: 1,
				},
			}
			Expect(util.GetCINamespacedName2("test", play).String()).To(Equal("p6e-cx-1/test-1"))
		})
		It("GetDeployNamespacedName", func() {
			By("build play")
			play := &ci.Play{
				Spec: ci.PlaySpec{
					Environment: "develop",
					ProjectID:   1,
				},
			}
			Expect(util.GetDeployNamespacedName("test", play).String()).To(Equal("test-develop-1/test-1"))
		})
		It("GetCILabels", func() {
			By("build play")
			play := &ci.Play{
				Spec: ci.PlaySpec{
					ProjectID:  1,
					PipelineID: 1,
				},
			}
			Expect(util.GetCILabels(play)).To(Equal(map[string]string{
				"projectid":  strconv.Itoa(int(play.Spec.ProjectID)),
				"pipelineid": strconv.Itoa(int(play.Spec.PipelineID)),
			}))
		})
		It("GetDockerImageTag", func() {
			var err error
			By("build play")
			play := &ci.Play{
				Spec: ci.PlaySpec{
					Name:      "test",
					ProjectID: 1,
					Commit: ci.Commit{
						SHA: "f2c2ab7eb8afb44ad94d3001b777f0d1dd35de33",
						Ref: "main",
					},
				},
			}

			By("no DockerURL")
			_, err = util.GetDockerImageTag(play)
			Expect(err).To(Succeed())

			By("Set DockerURL")
			play.Spec.DockerURL = "reg.example.io/test"
			_, err = util.GetDockerImageTag(play)
			Expect(err).To(Succeed())

			By("Set a wrong DockerURL")
			play.Spec.DockerURL = "http://{}"
			_, err = util.GetDockerImageTag(play)
			Expect(err).ToNot(Succeed())

		})
		It("GetDockerImageTagRaw", func() {
			var err error
			By("build play")
			play := &ci.Play{
				Spec: ci.PlaySpec{
					Name:      "test",
					ProjectID: 1,
					Commit: ci.Commit{
						SHA: "f2c2ab7eb8afb44ad94d3001b777f0d1dd35de33",
						Ref: "main",
					},
				},
			}

			By("no DockerURL")
			_, _, _, err = util.GetDockerImageTagRaw(play)
			Expect(err).To(Succeed())

			By("Set DockerURL")
			play.Spec.DockerURL = "reg.example.io/test"
			_, _, _, err = util.GetDockerImageTagRaw(play)
			Expect(err).To(Succeed())

			By("Set DockerURL without proto")
			play.Spec.DockerURL = "reg"
			_, _, _, err = util.GetDockerImageTagRaw(play)
			Expect(err).To(Succeed())

			By("Set DockerURL with tcp scheme")
			play.Spec.DockerURL = "tcp://reg:55/test:"
			_, _, _, err = util.GetDockerImageTagRaw(play)
			Expect(err).To(Succeed())

			By("Failed DockerURL with tcp scheme")
			play.Spec.DockerURL = "tcp://{}"
			_, _, _, err = util.GetDockerImageTagRaw(play)
			Expect(err).ToNot(Succeed())

		})
		It("IgnoreNotExists", func() {
			By("error is nil")
			Expect(util.IgnoreNotExists(nil)).To(BeNil())

			By("error raised an issue")
			Expect(util.IgnoreNotExists(errors.New("test")).Error()).To(Equal("test"))
		})
		It("GetObjectContain", func() {
			By("Set resource")
			play := &ci.Play{
				Spec: ci.PlaySpec{
					Name:      "test",
					ProjectID: 1,
					Commit: ci.Commit{
						SHA: "f2c2ab7eb8afb44ad94d3001b777f0d1dd35de33",
						Ref: "main",
					},
				},
			}
			Expect(util.GetObjectContain(play)).ToNot(Equal("<ERROR>\n"))
		})
		It("IsBuildStage", func() {
			By("build resource")
			play := &ci.Play{
				Spec: ci.PlaySpec{
					Stack: ci.Stack{
						Language: "",
					},
					Tasks: []map[ci.TaskType]ci.Task{
						{
							ci.Clean: ci.Task{},
						},
					},
				},
			}
			Expect(util.IsBuildStage(play)).To(Equal(false))

			By("Set build task")
			play.Spec.Tasks = []map[ci.TaskType]ci.Task{
				{
					ci.Build: ci.Task{},
				},
			}
			By("play language is ios")
			play.Spec.Stack.Language = "ios"
			Expect(util.IsBuildStage(play)).To(Equal(false))
		})
	})
})

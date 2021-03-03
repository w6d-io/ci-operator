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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/tekton/pipelinerun"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("integration test in pipeline run", func() {
	Context("setting", func() {
		It("does", func() {
			p := pipelinerun.PipelineRun{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						Spec: ci.PlaySpec{
							Tasks: []map[ci.TaskType]ci.Task{
								{
									ci.IntegrationTests: ci.Task{
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
			err := p.SetIntTest(0, ctrl.Log)
			Expect(err).To(Succeed())
		})
	})
})

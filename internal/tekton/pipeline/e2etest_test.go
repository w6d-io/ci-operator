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
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/tekton/pipeline"

	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Set e2e test", func() {
	Context("build", func() {
		It("succeed", func() {
			p := pipeline.Pipeline{
				Play: &ci.Play{
					Spec: ci.PlaySpec{
						ProjectID: 1,
					},
				},
			}
			err := p.SetPipelineE2ETest(ctrl.Log)
			Expect(err).To(Succeed())
		})
	})
})

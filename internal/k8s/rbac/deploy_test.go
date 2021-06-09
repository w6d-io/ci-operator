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
package rbac_test

import (
	"context"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/rbac"
	"github.com/w6d-io/ci-operator/internal/util"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RBAC", func() {
	Context("Check", func() {
		BeforeEach(func() {
			err := config.New("testdata/config.yaml")
			Expect(err).To(Succeed())
		})
		It("get role binding", func() {
			p := &ci.Play{
				Spec: ci.PlaySpec{
					ProjectID:  1,
					PipelineID: 1,
				},
			}
			r := rbac.GetRoleBinding(p)
			Expect(r).ToNot(BeNil())
		})
		It("get subject", func() {
			p := &ci.Play{
				Spec: ci.PlaySpec{
					ProjectID:  1,
					PipelineID: 1,
				},
			}
			r := rbac.GetSubject(p)
			Expect(r).ToNot(BeNil())
		})
	})
	When("in create method", func() {
		var (
			play *ci.Play
			d    rbac.Deploy
		)
		BeforeEach(func() {
			play = &ci.Play{
				Spec: ci.PlaySpec{
					Environment: "test",
					ProjectID:   1,
					PipelineID:  1,
				},
			}
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: util.GetDeployNamespacedName(config.GetDeployPrefix(), play).Namespace,
				},
			}
			_ = k8sClient.Create(context.TODO(), ns)
			err := config.New("testdata/config.yaml")
			Expect(err).To(Succeed())
			d = rbac.Deploy{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: play,
				},
			}
		})
		Context("we try", func() {
			It("succeed creation", func() {
				_ = d.Create(context.TODO(), k8sClient, ctrl.Log)
				d.Play.Spec.Tasks = []map[ci.TaskType]ci.Task{
					{
						"deploy": ci.Task{},
					},
				}
				err := d.Create(context.TODO(), k8sClient, ctrl.Log)
				Expect(err).To(Succeed())
				r := rbacv1.RoleBinding{}
				name := util.GetCINamespacedName2(rbac.Prefix, play).Name
				ns := util.GetDeployNamespacedName(config.GetDeployPrefix(), play).Namespace
				err = k8sClient.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: ns}, &r)
				Expect(err).To(Succeed())
				Expect(len(r.Subjects)).To(Equal(2))
			})
			It("with subject already exist", func() {
				d.Play.Spec.Tasks = []map[ci.TaskType]ci.Task{
					{
						ci.Deploy: ci.Task{},
					},
				}
				err := d.Create(context.TODO(), k8sClient, ctrl.Log)
				Expect(err).To(Succeed())
			})
			It("succeed update", func() {
				d.Play.Spec.Tasks = []map[ci.TaskType]ci.Task{
					{
						ci.Deploy: ci.Task{},
					},
				}
				play.Spec.PipelineID = 2
				err := d.Create(context.TODO(), k8sClient, ctrl.Log)
				Expect(err).To(Succeed())
				r := rbacv1.RoleBinding{}
				name := util.GetCINamespacedName2(rbac.Prefix, play).Name
				ns := util.GetDeployNamespacedName(config.GetDeployPrefix(), play).Namespace
				err = k8sClient.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: ns}, &r)
				Expect(err).To(Succeed())
				Expect(len(r.Subjects)).To(Equal(3))
			})
		})
	})
})

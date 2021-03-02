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
Created on 01/03/2021
*/
package task_test

import (
    "context"

    "github.com/w6d-io/ci-operator/internal/tekton/task"

    tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
    ci "github.com/w6d-io/ci-operator/api/v1alpha1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    ctrl "sigs.k8s.io/controller-runtime"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("e2e test", func() {
    Context("Build", func() {
        It("Failed to return step", func() {
            t := &task.Task{
                Index: 0,
                Play: &ci.Play{
                    Spec: ci.PlaySpec{
                        Tasks: []map[ci.TaskType]ci.Task{},
                    },
                },
                Client: k8sClient,
            }
            err := t.E2ETest(context.TODO(), ctrl.Log)
            Expect(err).ToNot(Succeed())
            Expect(err.Error()).To(Equal("no such task"))
        })
        It("not find step", func() {
            s := &ci.Step{
                ObjectMeta: metav1.ObjectMeta{
                    Name: "test-step-1",
                    Namespace: "default",
                    Annotations: map[string]string{
                        ci.AnnotationLanguage: "none",
                    },
                },
                Step: ci.StepSpec{

                },
            }
            err := k8sClient.Create(context.TODO(), s)
            Expect(err).To(Succeed())
            t := &task.Task{
                Index: 0,
                Play: &ci.Play{
                    Spec: ci.PlaySpec{
                        Stack: ci.Stack{
                            Language: "test",
                            Package: "test",
                        },
                        Tasks: []map[ci.TaskType]ci.Task{
                            {
                                ci.E2ETests: ci.Task{
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
                Client: k8sClient,
            }
            err = t.E2ETest(context.TODO(), ctrl.Log)
            Expect(err).ToNot(Succeed())
            Expect(err.Error()).To(ContainSubstring("no step found for e2e-tests"))

        })

        When("Step generic exists", func() {
            BeforeEach(func() {
                s := &ci.Step{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:      "test-step-2",
                        Namespace: "default",
                        Annotations: map[string]string{
                            ci.AnnotationKind:  "generic",
                            ci.AnnotationTask:  ci.E2ETests.String(),
                            ci.AnnotationOrder: "0",
                        },
                    },
                    Step: ci.StepSpec{
                        Step: tkn.Step{
                            Script: "echo test",
                        },
                    },
                }
                err := k8sClient.Create(context.TODO(), s)
                Expect(err).To(Succeed())
            })
            It("finds step", func() {
                t := &task.Task{
                    Index: 0,
                    Play: &ci.Play{
                        Spec: ci.PlaySpec{
                            Stack: ci.Stack{
                                Language: "test",
                                Package: "test",
                            },
                            Tasks: []map[ci.TaskType]ci.Task{
                                {
                                    ci.E2ETests: ci.Task{
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
                    Client: k8sClient,
                }
                var err error
                Eventually(func() bool {
                    err = t.E2ETest(context.TODO(), ctrl.Log)
                    return err == nil
                })
            })
        })
    })
    Context("Create", func() {
        It("", func() {
            e := &task.E2ETestTask{
                Meta: task.Meta{
                    Play: &ci.Play{
                        ObjectMeta: metav1.ObjectMeta{
                            Name: "test-create-1",
                            Namespace: "default",
                            UID: "uuid-uuid-uuid-uuid",
                        },
                        Spec: ci.PlaySpec{
                            PipelineID: 1,
                            ProjectID: 1,
                        },
                    },
                    Steps: []tkn.Step{
                        {
                            Script: `
echo "toto"
`,
                        },
                    },
                },
            }
            err := e.Create(context.TODO(), k8sClient, ctrl.Log)
            Expect(err).ToNot(Succeed())
            Expect(err.Error()).To(ContainSubstring("cross-namespace"))
        })
    })
})
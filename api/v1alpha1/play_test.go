/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 22/05/2021
*/

package v1alpha1_test

import (
	"github.com/w6d-io/ci-operator/api/v1alpha1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func() {
	Context("", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("", func() {
			By("build play")
			play := &v1alpha1.Play{
				Spec: v1alpha1.PlaySpec{
					Name: "test-project",
					Stack: v1alpha1.Stack{
						Language: "js",
						Package:  "npm",
					},
					Environment: "production",
					ProjectID:   1,
					PipelineID:  2,
					RepoURL:     "https://github.com/w6d-io/nodejs-sample",
					Commit: v1alpha1.Commit{
						SHA:     "1234567890",
						Ref:     "main",
						Message: "init",
					},
					Domain: "nodejs-sample.wildcard.sh",
					Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
						{
							v1alpha1.Build: v1alpha1.Task{
								Image: "docker",
								Script: []string{
									`echo "test1"`,
									`echo "test2"`,
								},
								Docker: v1alpha1.Docker{
									Filepath: "Dockerfile",
									Context:  ".",
								},
								Namespace: "default",
							},
						},
					},
					DockerURL: "docker.io/w6dio/nodejs-sample:latest",
					Secret: map[v1alpha1.SecretKind]string{
						v1alpha1.KubeConfig: "config",
					},
					Vault: &v1alpha1.Vault{
						Token: "token",
						Secrets: map[v1alpha1.SecretKind]v1alpha1.VaultSecret{
							v1alpha1.KubeConfig: {
								VolumePath: "/secret/key",
							},
						},
					},
				},
			}

			Expect(play.Get("play.name")).To(Equal("test-project"))
			Expect(play.Get("play.environment")).To(Equal("production"))
			Expect(play.Get("play.project_id")).To(Equal("1"))
			Expect(play.Get("play.pipeline_id")).To(Equal("2"))
			Expect(play.Get("play.repo_url")).To(Equal("https://github.com/w6d-io/nodejs-sample"))
			Expect(play.Get("play.commit.sha")).To(Equal("1234567890"))
			Expect(play.Get("play.commit.ref")).To(Equal("main"))
			Expect(play.Get("play.commit.message")).To(Equal("init"))
			Expect(play.Get("play.domain")).To(Equal("nodejs-sample.wildcard.sh"))
			Expect(play.Get("play.tasks.build.image")).To(Equal("docker"))
			Expect(play.Get("play.tasks.build.script")).To(Equal(`echo "test1"
echo "test2"`))
			Expect(play.Get("play.tasks.build.docker.filepath")).To(Equal("Dockerfile"))
			Expect(play.Get("play.tasks.build.docker.context")).To(Equal("."))
			Expect(play.Get("play.tasks.build.namespace")).To(Equal("default"))
			Expect(play.Get("play.docker_url")).To(Equal("docker.io/w6dio/nodejs-sample:latest"))
			Expect(play.Get("play.secret.kubeconfig")).To(Equal("play.secret.kubeconfig"))
			Expect(play.Get("play.vault.secrets.kubeconfig")).To(Equal("play.vault.secrets.kubeconfig"))
			Expect(play.Get("play")).To(Equal("play"))
			Expect(play.Get("bad.key")).To(Equal("bad.key"))
		})
	})
})

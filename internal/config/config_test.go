package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/w6d-io/ci-operator/internal/config"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Config", func() {
	Describe("Parse yaml configuration", func() {
// 		err := config.New("../../test/config/")
// 		if err != nil {
// 			Fail(err.Error())
// 		}
		Context("Manage config issues", func() {
			It("File does not exist", func() {
				Expect(config.New("../../test/config/no-file.yaml")).ToNot(BeNil())
			})
			It("load bad file", func() {
				Expect(config.New("../../test/config/file3.yaml")).ToNot(BeNil())
			})
			It("Webhook part misconfigured ", func() {
				Expect(config.New("../../test/config/file4.yaml")).ToNot(BeNil())
			})
		})
 		Context("Validate config", func() {
			It("File exists", func() {
				Expect(config.New("../../test/config/file1.yaml")).To(BeNil())
			})
 			It("Check mandatory value", func() {
 				Expect(config.Validate()).To(BeNil())
			})
			It("load file without mandatory part", func() {
				Expect(config.New("../../test/config/file2.yaml")).To(BeNil())
			})
			It("Missing mandatory part", func() {
				Expect(config.Validate()).ToNot(BeNil())
			})
			It("load tiny config ", func() {
				Expect(config.New("../../test/config/file5.yaml")).To(BeNil())
			})
			It("GetConfig function", func() {
				Expect(config.GetConfig()).To(Equal(&config.Config{
					DefaultDomain: "example.ci",
					Ingress: config.Ingress{
						Class: "nginx",
						Issuer: "letsencrypt-prod",
					},
					Volume: tkn.WorkspaceBinding{
						Name: "ws",
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
					DeployPrefix: "test",
				}))
			})
			It("GetConfigRaw function", func() {
				Expect(len(config.GetConfigRaw())).Should(Equal(10))
			})
			It("GetClusterRole function", func() {
				Expect(config.GetClusterRole()).Should(Equal(""))
			})
			It("GetDeployPrefix function", func() {
				Expect(config.GetDeployPrefix()).Should(Equal("test"))
			})
			It("GetNamespace function", func() {
				Expect(config.GetNamespace()).Should(Equal(""))
			})
			It("GetMinio function", func() {
				Expect(config.GetMinio()).To(Equal(&config.Minio{}))
			})
			It("GetMinioRaw function", func() {
				Expect(len(config.GetMinioRaw())).To(Equal(0))
			})
			It("GetRaw function", func() {
				Expect(config.GetRaw(make(chan int))).To(BeNil())
			})
			It("GetRaw function", func() {
				Expect(config.GetRaw("{test")).To(BeNil())
			})

 		})
		Context("Check Minio Method", func() {
			It("load tiny config ", func() {
				Expect(config.New("../../test/config/file5.yaml")).To(BeNil())
			})
			It("GetHost method", func() {
				Expect(config.GetMinio().GetHost()).To(Equal(""))
			})
			It("GetAccessKey method", func() {
				Expect(config.GetMinio().GetAccessKey()).To(Equal(""))
			})
			It("GetSecretKey method", func() {
				Expect(config.GetMinio().GetSecretKey()).To(Equal(""))
			})
			It("GetBucket method", func() {
				Expect(config.GetMinio().GetBucket()).To(Equal(""))
			})
		})
		Context("Test config function", func() {
			It("Load config", func() {
				Expect(config.New("../../test/config/file1.yaml")).To(BeNil())
			})
			It("Workspace function", func() {
				Expect(len(config.Workspaces())).Should(Equal(4))
			})
			It("Volume function", func() {
				Expect(config.Volume().VolumeClaimTemplate.Spec.Resources.Requests.Storage().String()).
					To(Equal("2Gi"))
			})
			It("GetWorkspacePath function", func() {
				Expect(config.GetWorkspacePath("artifacts")).To(Equal("/artifacts"))
			})
			It("GetWorkspacePath function with unset value", func() {
				Expect(config.GetWorkspacePath("test")).To(Equal(""))
			})

		})
	})
})

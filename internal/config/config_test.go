package config_test

import (
	"github.com/w6d-io/ci-operator/internal/config"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("Parse yaml configuration", func() {
		Context("Manage config issues", func() {
			It("File does not exist", func() {
				Expect(config.New("testdata/no-file.yaml")).ToNot(BeNil())
			})
			It("load bad file", func() {
				Expect(config.New("testdata/file3.yaml")).ToNot(BeNil())
			})
			It("Webhook part misconfigured ", func() {
				Expect(config.New("testdata/file4.yaml")).ToNot(BeNil())
			})
		})
		Context("Validate config", func() {
			It("File exists", func() {
				Expect(config.New("testdata/file1.yaml")).To(BeNil())
			})
			It("Check mandatory value", func() {
				Expect(config.Validate()).To(BeNil())
			})
			It("load file without mandatory part", func() {
				Expect(config.New("testdata/file2.yaml")).To(BeNil())
			})
			It("Missing mandatory part", func() {
				Expect(config.Validate()).ToNot(BeNil())
			})
			It("load tiny config ", func() {
				Expect(config.New("testdata/file5.yaml")).To(BeNil())
			})
			It("GetConfig function", func() {
				Expect(config.GetConfig()).To(Equal(&config.Config{
					DefaultDomain: "example.ci",
					Ingress: config.Ingress{
						Class:  "nginx",
						Issuer: "letsencrypt-prod",
					},
					Volume: tkn.WorkspaceBinding{
						Name:     "ws",
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
					DeployPrefix: "test",
				}))
			})
			It("GetConfigRaw function", func() {
				Expect(len(config.GetConfigRaw())).Should(Equal(11))
			})
			It("GetClusterRole function", func() {
				Expect(config.GetClusterRole()).Should(Equal(""))
			})
			It("GetDeployPrefix function", func() {
				Expect(config.GetDeployPrefix()).Should(Equal("test"))
			})
			It("GetNamespace function", func() {
				By("set namespace")
				config.SetNamespace("default")

				Expect(config.GetNamespace()).Should(Equal("default"))
			})
			It("GetMinio function", func() {
				By("minio isn't configured")
				Expect(config.GetMinio()).To(Equal(&config.Minio{}))

				By("load a config")
				err := config.New("testdata/file1.yaml")
				Expect(err).To(Succeed())

				By("minio is configured")
				Expect(config.GetMinio()).ToNot(BeNil())
				Expect(config.GetMinio().GetBucket()).To(Equal("values"))

			})
			It("GetMinioRaw function", func() {
				By("no minio entry")
				Expect(config.New("testdata/file5.yaml")).To(Succeed())
				Expect(len(config.GetMinioRaw())).To(Equal(0))

				By("no minio entry")
				Expect(config.New("testdata/file1.yaml")).To(Succeed())
				Expect(len(config.GetMinioRaw())).To(Equal(4))
				//config.GetMinioRaw()
			})
			It("GetRaw function", func() {
				Expect(config.GetRaw(make(chan int))).To(BeNil())
			})
			It("GetRaw function", func() {
				Expect(config.GetRaw("{test")).To(BeNil())
			})
			It("gets values", func() {
				Expect(config.GetValues()).ToNot(BeNil())
			})
			It("gets hash", func() {
				By("no hash entry")
				Expect(config.New("testdata/file5.yaml")).To(Succeed())
				Expect(config.GetHash()).ToNot(BeNil())

				Expect(config.New("testdata/file1.yaml")).To(Succeed())
				Expect(config.GetHash()).ToNot(BeNil())
				Expect(config.GetHash().GetSalt()).To(Equal("wildcard"))
				Expect(config.GetHash().GetMinLength()).To(Equal(16))
			})
			It("gets vault", func() {
				Expect(config.New("testdata/file1.yaml")).To(Succeed())
				Expect(config.GetVault()).ToNot(BeNil())
				Expect(config.GetVault().GetHost()).To(Equal("vault.svc:8200"))
				Expect(config.GetVault().GetToken()).To(Equal("token"))
			})

		})
		Context("Check Minio Method", func() {
			It("load tiny config ", func() {
				By("load config")
				Expect(config.New("testdata/file5.yaml")).To(Succeed())
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
				Expect(config.New("testdata/file1.yaml")).To(BeNil())
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
			It("GetWebhooks function", func() {
				Expect(len(config.GetWebhooks())).To(Equal(1))
			})
			It("gets pod template", func() {
				By("when pod template does not exist")
				Expect(config.PodTemplate()).ToNot(BeNil())

				By("when pod template does exist")
				Expect(config.PodTemplate()).ToNot(BeNil())
			})
			It(" ", func() {
			})
		})
	})
})

package config_test

import (
	"github.com/w6d-io/ci-operator/internal/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

)

var _ = Describe("Config", func() {
	Describe("Parse yaml configuration", func() {
// 		err := config.New("../../test/config/")
// 		if err != nil {
// 			Fail(err.Error())
// 		}
		Context("Load config file", func() {
			It("File does not exist", func() {
				Expect(config.New("../../test/config/no-file.yaml")).ToNot(BeNil())
			})
			It("File exists", func() {
				Expect(config.New("../../test/config/file1.yaml")).To(BeNil())
			})
		})
 		Context("Check element in config", func() {
 			It("Validate")
 			It("Default domain", func() {
 				Expect(config.GetConfig().DefaultDomain).To(Equal("example.ci"))
 			})
 		})
	})
})

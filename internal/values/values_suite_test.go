/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 15/02/2021
*/

package values_test

import (
	"github.com/w6d-io/ci-operator/internal/config"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestValues(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Values Suite")
}

var _ = BeforeSuite(func(done Done) {
	err := config.New("testdata/config.yaml")
	Expect(err).To(Succeed())
	close(done)
}, 60)

var _ = AfterSuite(func() {
})

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
Created on 21/12/2020
*/

package values_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/w6d-io/ci-operator/internal/config"
	"testing"

	"github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/values"

	f "k8s.io/apimachinery/pkg/fields"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Templates", func() {
	Context("", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("returns empty because vault config isn't set", func() {
			var err error
			err = config.New("testdata/config1.yaml")
			Expect(err).To(Succeed())
			Expect(values.Vault("token", "path", "key")).To(Equal(""))
		})
		It("return empty because vault is unreachable", func() {
			var err error
			err = config.New("testdata/config.yaml")
			Expect(err).To(Succeed())
			Expect(values.Vault("token", "path", "key")).To(Equal(""))
		})
	})
})

func TestTemplates_PrintTemplate(t *testing.T) {
	type fields struct {
		Play     *v1alpha1.Play
		Spec     v1alpha1.PlaySpec
		Internal map[string]interface{}
	}
	type args struct {
		out   *bytes.Buffer
		name  string
		templ string
	}
	buf := new(bytes.Buffer)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test1",
			fields{
				&v1alpha1.Play{},
				v1alpha1.PlaySpec{
					Name: "test",
					Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
						{
							v1alpha1.Deploy: v1alpha1.Task{
								Variables: f.Set{
									"HOSTNAME": "w6d.io",
									"USERNAME": "mea",
								},
							},
						},
					},
					Domain: "mea-test.w6d.io",
				},
				map[string]interface{}{
					"domain": "example.ci",
					"hash": map[string]interface{}{
						"salt":       "wildcard",
						"min_length": 16,
					},
				},
			},
			args{
				buf,
				"test.yaml",
				values.HelmValuesTemplate,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &values.Templates{
				Client:   k8sClient,
				Play:     tt.fields.Play,
				Spec:     tt.fields.Spec,
				Internal: tt.fields.Internal,
			}
			if err := in.PrintTemplate(context.TODO(), tt.args.out, tt.args.name, tt.args.templ); (err != nil) != tt.wantErr {
				fmt.Printf("out = %s\n", tt.args.out.String())
				t.Errorf("PrintTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashID(t *testing.T) {
	type args struct {
		pid float64
	}
	tests := []struct {
		name     string
		previous func()
		post     func()
		args     args
		want     string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			"long id",
			func() {},
			func() {},
			args{
				1236541651684135186185431085485413.21,
			},
			"",
			true,
		},
		{
			"change_alphabet",
			func() {
				values.Alphabet = "abc"
			},
			func() {
				values.Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
			},
			args{
				0,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.previous()
			got, err := values.HashID(tt.args.pid)
			tt.post()
			if (err != nil) != tt.wantErr {
				t.Errorf("HashID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HashID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToYaml(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			"good yaml",
			args{
				`---
test: ok
`,
			},
			`|
    ---
    test: ok
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := values.ToYaml(tt.args.v); got != tt.want {
				t.Errorf("ToYaml() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

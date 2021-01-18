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

package values

import (
	"bytes"
	"fmt"
	"testing"

	f "k8s.io/apimachinery/pkg/fields"

	"github.com/w6d-io/ci-operator/api/v1alpha1"
)

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
					Dependencies: map[v1alpha1.DependencyType]v1alpha1.Dependency{
						"mongodb": {
							Variables: f.Set{
								"HOST":     "$DB_HOST",
								"PASSWORD": "$DB_PASSWORD",
								"USERNAME": "$DB_USERNAME",
							},
						},
					},
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
				HelmValuesTemplate,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &Templates{
				Play:     tt.fields.Play,
				Spec:     tt.fields.Spec,
				Internal: tt.fields.Internal,
			}
			if err := in.PrintTemplate(tt.args.out, tt.args.name, tt.args.templ); (err != nil) != tt.wantErr {
				fmt.Printf("out = %s\n", tt.args.out.String())
				t.Errorf("PrintTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

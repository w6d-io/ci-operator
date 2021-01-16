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
Created on 31/12/2020
*/

package util

import (
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"testing"
)

func TestIsBuildStage(t *testing.T) {
	type args struct {
		play *ci.Play
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"test1",
			args{
				play: &ci.Play{
					Spec: ci.PlaySpec{
						Tasks: []map[ci.TaskType]ci.Task{
							{
								ci.Build: ci.Task{},
							},
						},
					},
				},
			},
			true,
		},
		{"test2",
			args{
				play: &ci.Play{
					Spec: ci.PlaySpec{
						Tasks: []map[ci.TaskType]ci.Task{
							{
								ci.Clean: ci.Task{},
							},
						},
					},
				},
			},
			false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBuildStage(tt.args.play); got != tt.want {
				t.Errorf("IsBuildStage() = %v, want %v", got, tt.want)
			}
		})
	}
}

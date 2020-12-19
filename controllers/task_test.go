/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 23/11/2020
*/

package controllers

import (
	"github.com/w6d-io/ci-operator/api/v1alpha1"
	"testing"
)

func TestIsBuildStage(t *testing.T) {
	type args struct {
		play v1alpha1.Play
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"test1",
			args{
				play: v1alpha1.Play{
					Spec: v1alpha1.PlaySpec{
						Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
							{
								v1alpha1.Build: v1alpha1.Task{},
							},
						},
					},
				},
			},
			true,
		},
		{"test2",
			args{
				play: v1alpha1.Play{
					Spec: v1alpha1.PlaySpec{
						Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
							{
								v1alpha1.Clean: v1alpha1.Task{},
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

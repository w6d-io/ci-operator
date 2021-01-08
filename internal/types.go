/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 28/12/2020
*/

package internal

import (
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
)

type WorkFlowStruct struct {
	Scheme *runtime.Scheme
	Play   *ci.Play
}

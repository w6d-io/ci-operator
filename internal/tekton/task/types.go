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
Created on 07/01/2021
*/

package task

import (
	"context"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Task struct {
	// Index is the position of the task in the list
	Index   int
	Creates []func(context.Context, client.Client, logr.Logger) error
	Client  client.Client
	Play    *ci.Play
	Scheme  *runtime.Scheme
	Params  map[string][]ci.ParamSpec
}

// Meta is the base struct for all struct with create method
type Meta struct {
	Steps    []tkn.Step
	Sidecars []tkn.Sidecar
	Play     *ci.Play
	Scheme   *runtime.Scheme
}

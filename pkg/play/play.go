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
Created on 22/11/2020
*/

package play

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/util"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

// CreateCI takes a Play struct to create tekton Pipeline
func CreateCI(ctx context.Context, p *ci.Play, logger logr.Logger,
	r client.Client, scheme *runtime.Scheme) error {

	log := logger.WithName("CreateCI").WithValues("cx-namespace", util.InNamespace(p))
	p.Status.State = ci.Creating
	if err := r.Status().Update(ctx, p); err != nil {
		return err
	}
	var wf WFInterface
	wf = &WFType{
		Client: r,
		Scheme: scheme,
		Creates: []CIFunc{},
	}
//	wf = New(r, scheme)

	if err := wf.CreateValues(p, logger); err != nil {
		return err
	}
	if err := wf.ServiceAccount(p, logger); err != nil {
		return err
	}
	if err := wf.Rbac(p, logger); err != nil {
		return err
	}
	if err := wf.GitSecret(p, logger); err != nil {
		return err
	}
	if err := wf.DockerCredSecret(p, logger); err != nil {
		return err
	}
	if err := wf.MinIOSecret(p, logger); err != nil {
		return err
	}
	if err := wf.SetGitCreate(p, logger); err != nil {
		return err
	}
	if err := wf.SetImageCreate(p, logger); err != nil {
		return err
	}
	if err := wf.SetTask(ctx, p, logger); err != nil {
		return err
	}
	if err := wf.SetPipeline(p, logger); err != nil {
		return err
	}
	if err := wf.SetPipelineRun(p, logger); err != nil {
		return err
	}
	log.Info("Launch creation")
	if err := wf.Run(ctx, r, log); err != nil {
		log.Error(err, "CI creation failed")
		// TODO add rollback ( delete resource created before )
		return err
	}
	return nil
}

// TODO Add is ci exists function / method

//// New creates a WFInterface instance
//func New(client client.Client, scheme *runtime.Scheme) *WFType {
//	wf := new(WFType)
//	wf.Creates = []CIFunc{}
//	wf.Client = client
//	wf.Scheme = scheme
//	return wf
//}

// Run executes Create methods in WFType
func (wf *WFType) Run(ctx context.Context, r client.Client, log logr.Logger) error {
	for _, c := range wf.Creates {
		if err := c(ctx, r, log); err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

// Add appends the CIFunc to the run list
func (wf *WFType) Add(ciFunc CIFunc) error {
	wf.Creates = append(wf.Creates, ciFunc)
	return nil
}

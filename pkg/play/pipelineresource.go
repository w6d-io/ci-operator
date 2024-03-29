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
	"errors"
	"net/url"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/tekton/pipelineresource"
	"github.com/w6d-io/ci-operator/internal/util"
)

// SetGitCreate adds git pipeline resource
func (wf *WFType) SetGitCreate(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("SetGitCreate")
	log.V(1).Info("Check repository URL")
	URL, err := url.Parse(p.Spec.RepoURL)
	if err != nil {
		log.Error(err, "URL parse error")
		return err
	}
	git := &pipelineresource.GitPR{
		// TODO put the prefix in config
		NamespacedName: util.GetCINamespacedName("pr-git", p),
		URL:            URL,
		Labels:         util.GetCILabels(p),
		Revision:       p.Spec.Commit.SHA,
		Play:           p,
		Scheme:         wf.Scheme,
	}
	log.V(1).Info("add pipeline resource create method for git")
	if err := wf.Add(git.Create); err != nil {
		return err
	}
	// TODO add git token secret creation
	log.V(1).Info("GitCreate added")
	return nil
}

// SetImageCreate adds image pipeline resource
func (wf *WFType) SetImageCreate(p *ci.Play, logger logr.Logger) error {
	if util.IsBuildStage(p) {
		log := logger.WithName("SetImageCreate")
		URL, err := util.GetDockerImageTag(p)
		if err != nil {
			log.Error(err, "get repository address failed")
			return errors.New("get repository address failed")
		}
		image := pipelineresource.ImagePR{
			// TODO set the prefix in configuration
			NamespacedName: util.GetCINamespacedName("pr-image", p),
			// TODO allow
			URL:    URL,
			Labels: util.GetCILabels(p),
			Play:   p,
			Scheme: wf.Scheme,
		}
		log.V(1).Info("add pipeline resource create method for image")
		if err := wf.Add(image.Create); err != nil {
			return err
		}
		// TODO add registry credential secret creation
		log.V(1).Info("ImageCreate added")

	}
	return nil
}

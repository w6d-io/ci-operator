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
Created on 08/01/2021
*/

package pipelineresource

import (
	"context"
	"fmt"
	"net/url"
	"time"

	resourcev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/resource/v1alpha1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// GitPR git PipelineResource type for CI
type GitPR struct {
	NamespacedName types.NamespacedName
	Labels         map[string]string
	Revision       string
	URL            *url.URL
	Play           *ci.Play
	Scheme         *runtime.Scheme
}

// Create implements the CIInterface method
func (g GitPR) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("type", "pipelineResource", "for", "git")
	log.V(1).Info("creating")

	gitResource := &resourcev1alpha1.PipelineResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:        g.NamespacedName.Name,
			Namespace:   g.NamespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      g.Labels,
		},
		Spec: resourcev1alpha1.PipelineResourceSpec{
			Type: resourcev1alpha1.PipelineResourceTypeGit,
			Params: []resourcev1alpha1.ResourceParam{
				{
					Name:  "revision",
					Value: g.Revision,
				},
				{
					Name:  "URL",
					Value: g.URL.String(),
				},
			},
		},
	}

	// set the current time in the new pipeline resource git type resource in annotation
	gitResource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(g.Play, gitResource, g.Scheme); err != nil {
		return err
	}
	log.V(1).Info(fmt.Sprintf("pipelineResource contains\n%v",
		util.GetObjectContain(gitResource)))
	if err := r.Create(ctx, gitResource); err != nil {
		return err
	}
	return nil
}

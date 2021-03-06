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
Created on 28/12/2020
*/

package secrets

import (
	"context"
	"fmt"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"net/url"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/sa"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	// GitTektonAnnotation
	// TODO use an increment in order to be able to manage several repository git
	GitTektonAnnotation = "tekton.dev/git-0"
	GitPrefixSecret     = "secret-git"
	GitUsername         = "oauth2"
)

// GitCreate creates the git credential secret and add it into the service account
func (s *Secret) GitCreate(ctx context.Context, r client.Client, logger logr.Logger) error {
	log := logger.WithName("Create").WithValues("action", GitPrefixSecret)
	log.V(1).Info("creating")

	namespacedName := util.GetCINamespacedName(GitPrefixSecret, s.Play)
	domain, err := getHTTPDomain(s.Play.Spec.RepoURL, log)
	if err != nil {
		return err
	}
	resource := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(s.Play),
		},
		StringData: map[string]string{
			"username": GitUsername,
			"password": s.GetSecret(ci.GitToken, log),
		},
		Type: "kubernetes.io/basic-auth",
	}
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	resource.Annotations[GitTektonAnnotation] = domain
	if err := controllerutil.SetControllerReference(s.Play, resource, s.Scheme); err != nil {
		return err
	}
	log.V(1).Info(resource.Kind, "content", fmt.Sprintf("%v",
		util.GetObjectContain(resource)))
	// Create Secret
	if err := r.Create(ctx, resource); err != nil {
		log.Error(err, "create")
		return err
	}
	if err := sa.Update(ctx, resource.Name,
		util.GetCINamespacedName(sa.Prefix, s.Play), r); err != nil {
		log.Error(err, "update")
		return err
	}
	// All went well
	return nil
}

// getHTTPDomain returns the base url of the repository
func getHTTPDomain(repoURL string, log logr.Logger) (domain string, err error) {
	log = log.WithName("getHttpDomain")
	URL := new(url.URL)
	URL, err = url.Parse(repoURL)
	if err != nil {
		log.Error(err, "url parse returns")
		return "", err
	}
	return fmt.Sprintf("%s://%s", URL.Scheme, URL.Hostname()), nil
}

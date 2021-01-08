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

package secrets

import (
	"context"
	"fmt"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/serviceaccount"
	"github.com/w6d-io/ci-operator/internal/util"
	"net/url"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	// TODO use an increment in order to be able to manage several repository git
	GitTektonAnnotation = "tekton.dev/git-0"
	GitPrefixSecret     = "secret-git"
	GitSecretKey        = "git_token"
	GitUsername         = "oauth2"
)

// GitCreate creates the git credential secret and add it into the service account
func (s *Secret) GitCreate(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", GitPrefixSecret)
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
			"password": s.Play.Spec.Secret[GitSecretKey],
		},
		Type: "kubernetes.io/basic-auth",
	}
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	resource.Annotations[GitTektonAnnotation] = domain
	if err := controllerutil.SetControllerReference(s.Play, resource, s.Scheme); err != nil {
		return err
	}
	log.V(1).Info(fmt.Sprintf("Secret contains\n%v",
		util.GetObjectContain(resource)))
	// Create Secret
	if err := r.Create(ctx, resource); err != nil {
		return err
	}

	if err := serviceaccount.Update(ctx, resource.Name,
		util.GetCINamespacedName(serviceaccount.Prefix, s.Play), r); err != nil {
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

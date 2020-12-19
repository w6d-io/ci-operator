/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
*/

package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var playlog = logf.Log.WithName("play-resource")

func (in *Play) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-ci-w6d-io-v1alpha1-play,mutating=true,failurePolicy=fail,admissionReviewVersions=v1alpha1,sideEffects=None,groups=ci.w6d.io,resources=plays,verbs=create;update,versions=v1alpha1,name=mplay.kb.io

var _ webhook.Defaulter = &Play{}

func (in *Play) Default() {
	playlog.Info("default", "name", in.Name)

	if in.Spec.Scope.Name == "" {
		in.Spec.Scope.Name = "default"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-ci-w6d-io-v1alpha1-play,mutating=false,failurePolicy=fail,admissionReviewVersions=v1alpha1,sideEffects=None,groups=ci.w6d.io,resources=plays,versions=v1alpha1,name=vplay.kb.io

var _ webhook.Validator = &Play{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *Play) ValidateCreate() error {
	playlog.Info("validate create", "name", in.Name)

	// TODO(user): fill in your validation logic upon object creation.
	var allErrs field.ErrorList
	allErrs = in.validateTaskType()
	if len(allErrs) > 0 {
		return apierrors.NewInvalid(
			PlayGroupKind, in.Name, allErrs)
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *Play) ValidateUpdate(old runtime.Object) error {
	playlog.Info("validate update", "name", in.Name)

	// TODO(user): fill in your validation logic upon object update.
	var allErrs field.ErrorList
	allErrs = in.validateTaskType()
	if len(allErrs) > 0 {
		return apierrors.NewInvalid(
			PlayGroupKind, in.Name, allErrs)
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *Play) ValidateDelete() error {
	playlog.Info("validate delete", "name", in.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (in Play) validateTaskType() field.ErrorList {
	var taskErrs field.ErrorList
	for _, task := range in.Spec.Tasks {
		for t := range task {
			switch t {
			case Build,Sonar,UnitTests,IntegrationTests,Deploy,Clean:
				continue
			default:
				taskErrs = append(taskErrs,
					field.Invalid(field.NewPath("spec").Child("tasks"),
					t,
					"not a TaskType"))

			}
		}
	}
	return taskErrs
}


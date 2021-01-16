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

// +kubebuilder:webhook:path=/mutate-ci-w6d-io-v1alpha1-play,mutating=true,failurePolicy=fail,admissionReviewVersions=v1;v1beta1,sideEffects=None,groups=ci.w6d.io,resources=plays,verbs=create;update,versions=v1alpha1,name=mplay.kb.io

var _ webhook.Defaulter = &Play{}

func (in *Play) Default() {
	playlog.Info("default", "name", in.Name)

	if in.Spec.Scope.Name == "" {
		in.Spec.Scope.Name = "default"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-ci-w6d-io-v1alpha1-play,mutating=false,failurePolicy=fail,admissionReviewVersions=v1;v1beta1,sideEffects=None,groups=ci.w6d.io,resources=plays,versions=v1alpha1,name=vplay.kb.io

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

	var allErrs field.ErrorList
	allErrs = in.validateTaskType()

	if old.(*Play).Spec.PipelineID != in.Spec.PipelineID {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("pipelineID"),
				in.Spec.PipelineID,
				"pipelineID cannot be changed"))
	}
	if old.(*Play).Spec.ProjectID != in.Spec.ProjectID {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("projectID"),
				in.Spec.ProjectID,
				"pipelineID cannot be changed"))
	}
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
			case Build, Sonar, UnitTests, IntegrationTests, Deploy, Clean:
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

/*
Copyright 2021.

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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"net/url"
	"reflect"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// log is for logging in this package.
var playlog = logf.Log.WithName("play-resource")

func (in *Play) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-ci-w6d-io-v1alpha1-play,mutating=true,failurePolicy=fail,admissionReviewVersions=v1;v1beta1,sideEffects=None,groups=ci.w6d.io,resources=plays,verbs=create;update,versions=v1alpha1,name=mutate.play.ci.w6d.io

var _ webhook.Defaulter = &Play{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *Play) Default() {
	playlog.Info("default", "name", in.Name)

	if in.Spec.Scope.Name == "" {
		in.Spec.Scope.Name = "default"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-ci-w6d-io-v1alpha1-play,mutating=false,failurePolicy=fail,admissionReviewVersions=v1;v1beta1,sideEffects=None,groups=ci.w6d.io,resources=plays,versions=v1alpha1,name=validate.play.ci.w6d.io

var _ webhook.Validator = &Play{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *Play) ValidateCreate() error {
	playlog.Info("validate create", "name", in.Name)
	var allErrs field.ErrorList
	allErrs = in.validateTaskType()
	allErrs = append(allErrs, in.commonValidation()...)
	allErrs = append(allErrs, in.validateVault()...)
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(
		schema.GroupKind{Group: "ci.w6d.io", Kind: "Play"},
		in.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *Play) ValidateUpdate(old runtime.Object) error {
	playlog.Info("validate update", "name", in.Name)

	var allErrs field.ErrorList
	allErrs = in.validateTaskType()
	allErrs = append(allErrs, in.commonValidation()...)
	allErrs = append(allErrs, in.validateVault()...)
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
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(
		schema.GroupKind{Group: "ci.w6d.io", Kind: "Play"},
		in.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *Play) ValidateDelete() error {
	playlog.Info("validate delete", "name", in.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (in Play) validateTaskType() field.ErrorList {
	var taskErrs field.ErrorList
	if len(in.Spec.Tasks) == 0 {
		taskErrs = append(taskErrs,
			field.Invalid(field.NewPath("spec").Child("tasks"),
				in.Spec.Tasks,
				"tasks cannot be empty"))
	}
	for _, task := range in.Spec.Tasks {
		for t := range task {
			switch t {
			case Build, Sonar, UnitTests, IntegrationTests, Deploy, Clean, E2ETests:
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

func (in *Play) commonValidation() field.ErrorList {
	playlog.Info("validate common", "name", in.Name)
	var allErrs field.ErrorList
	if in.Spec.ProjectID == 0 {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("project_id"),
				in.Spec.ProjectID,
				"cannot be 0"))
	}
	if in.Spec.PipelineID == 0 {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("pipeline_id"),
				in.Spec.PipelineID,
				"cannot be 0"))
	}
	if in.Spec.Environment == "" {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("environment"),
				in.Spec.Environment,
				"environment cannot be empty"))
	}
	if in.Spec.Commit.Ref == "" {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("commit").Child("ref"),
				in.Spec.Commit.Ref,
				"cannot be empty"))
	}
	if in.Spec.RepoURL == "" {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("repo_url"),
				in.Spec.RepoURL,
				"repo_url cannot be empty"))
	} else {
		if _, err := url.Parse(in.Spec.RepoURL); err != nil {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec").Child("repo_url"),
					in.Spec.RepoURL,
					err.Error()))
		}
	}
	if in.Spec.Commit.SHA == "" {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("commit").Child("sha"),
				in.Spec.Commit.SHA,
				"cannot be empty"))
	}
	if in.Spec.Domain != "" {
		if !validateDomain(in.Spec.Domain) {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec").Child("domain"),
					in.Spec.Domain,
					"domain invalid"))
		}
	}
	return allErrs
}

func validateDomain(domain string) bool {
	pattern := `^([a-z0-9]{1}[a-z0-9\-]{0,62}){1}(\.[a-z0-9]{1}[a-z0-9\-]{0,62})*[\._]?$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(domain)
}

func (in *Play) validateVault() field.ErrorList {
	var allErrs field.ErrorList

	if in.Spec.Vault != nil {
		for secret, _ := range in.Spec.Vault.Secrets {
			if ok, _ := inArray(secret, SecretKinds); !ok {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec").Child("vault").Child("secrets"),
						secret,
						"secret kind not supported"))
			}
		}
	}
	return allErrs
}

func inArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

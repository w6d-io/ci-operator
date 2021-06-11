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

package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/w6d-io/ci-operator/internal/tekton/pipelinerun"
	"github.com/w6d-io/ci-operator/internal/util"
	"github.com/w6d-io/ci-operator/pkg/play"
	"github.com/w6d-io/ci-operator/pkg/webhook"
	"github.com/w6d-io/hook"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

// PlayReconciler reconciles a Play object
type PlayReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ci.w6d.io,resources=plays,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ci.w6d.io,resources=plays/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ci.w6d.io,resources=plays/finalizers,verbs=update
//+kubebuilder:rbac:groups=ci.w6d.io,resources=limitcis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ci.w6d.io,resources=limitcis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ci.w6d.io,resources=limitcis/finalizers,verbs=update
//+kubebuilder:rbac:groups=ci.w6d.io,resources=steps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ci.w6d.io,resources=steps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ci.w6d.io,resources=steps/finalizers,verbs=update
//+kubebuilder:rbac:groups=tekton.dev,resources=pipelineresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tekton.dev,resources=pipelineresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tekton.dev,resources=pipelineruns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tekton.dev,resources=pipelineruns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tekton.dev,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tekton.dev,resources=pipelines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tekton.dev,resources=taskruns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tekton.dev,resources=taskruns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tekton.dev,resources=tasks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tekton.dev,resources=tasks/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PlayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	correlationID := uuid.New().String()
	ctx = util.NewCorrelationIDContext(ctx, correlationID)
	ctx = util.NewPlayContext(ctx, req.NamespacedName.String())
	logger := r.Log.WithValues("play", req.NamespacedName, "correlation_id", correlationID)
	log := logger.WithName("Reconcile")
	// get the play resource
	p := new(ci.Play)

	if err := r.Get(ctx, req.NamespacedName, p); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Play resource not found. Ignore since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch Play")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	p.Spec.Name = strings.ToLower(p.Spec.Name)
	p.Spec.Environment = strings.ToLower(p.Spec.Environment)
	log.V(1).Info("get", "pipelinerun", util.GetCINamespacedName(pipelinerun.Prefix, p).String())
	var childPr tkn.PipelineRun
	err := r.Get(ctx, util.GetCINamespacedName(pipelinerun.Prefix, p), &childPr)
	if client.IgnoreNotFound(err) != nil {
		log.Error(err, "Unable to get PipelineRun")
		return ctrl.Result{}, err
	}

	if !apierrors.IsNotFound(err) {
		if childPr.Name != "" && p.Status.PipelineRunName != childPr.Name {
			if err := r.UpdateName(ctx, p, childPr.Name); err != nil {
				return ctrl.Result{Requeue: true}, client.IgnoreNotFound(err)
			}
		}
		if err := r.UpdateStatus(ctx, p, util.Condition(childPr.Status.Conditions),
			util.Message(childPr.Status.Conditions)); err != nil {
			log.Error(err, "unable to update Play status")
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{Requeue: false}, nil
	}
	log.V(1).Info("pipelinerun not found")
	log.V(1).Info("getting all pipeline run")
	var prs tkn.PipelineRunList
	if err := r.List(ctx, &prs, util.InNamespace(p)); util.IgnoreNotExists(err) != nil {
		log.Error(err, "Unable to list PipelineRuns in ", util.InNamespace(p))
		log.V(1).Info("update status", "status", ci.Errored,
			"step", "2")
		if err := r.UpdateStatus(ctx, p, ci.Errored, err.Error()); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	}

	log.V(1).Info("check pipeline run running")
	var runningPipeline []tkn.PipelineRun
	for _, pr := range prs.Items {
		if util.IsPipelineRunning(pr) {
			runningPipeline = append(runningPipeline, pr)
		}
	}
	log.V(1).Info("pipelinerun", "running", len(runningPipeline))

	log.V(1).Info("get limitCi")
	var limits ci.LimitCiList
	if err := r.List(ctx, &limits, util.InNamespace(p)); client.IgnoreNotFound(err) != nil {
		log.Error(err, "unable to list LimitCi")
		log.V(1).Info("update status", "status", ci.Errored,
			"step", "3")
		if err := r.UpdateStatus(ctx, p, ci.Errored, err.Error()); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{Requeue: true}, err
	}

	if len(limits.Items) > 0 && (limits.Items[0].Spec.Concurrent <= int64(len(runningPipeline))) {
		log.V(1).Info("limit ci", "action", "queued")
		log.V(1).Info("update status", "status", ci.Queued,
			"step", "4")
		if err := r.UpdateStatus(ctx, p, ci.Queued, ""); err != nil {
			return ctrl.Result{}, err
		}
		lp := webhook.GetLimitPayload(p, limits.Items[0], "concurrent")
		if err := hook.Send(lp, ctrl.Log, "concurrent"); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
		//return ctrl.Result{}, nil
	}
	err = play.CreateCI(ctx, p, logger, r.Client, r.Scheme)
	if err != nil {
		log.Error(err, "Failed to create CI")
		log.V(1).Info("update status", "status", ci.Errored,
			"step", "5")
		if err := r.UpdateStatus(ctx, p, ci.Errored, err.Error()); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{Requeue: true}, err
	}

	payload := webhook.GetPayLoad(p)
	log.V(1).Info("hook process")
	if err := hook.Send(payload, ctrl.Log, "end"); err != nil {
		return ctrl.Result{Requeue: false}, err
	}

	log.V(1).Info("update status", "status", "---",
		"step", "6")
	if err := r.UpdateStatus(ctx, p, "---", ""); err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{Requeue: false}, nil
}

// UpdateStatus set the status of tekton resource state
func (r *PlayReconciler) UpdateStatus(ctx context.Context, p *ci.Play, state ci.State, message string) error {
	correlationID, _ := util.GetCorrelationIDFromContext(ctx)
	nn, _ := util.GetPlayFromContext(ctx)
	log := ctrl.Log.WithName("Reconcile").WithName("UpdateStatus").WithValues("correlation_id", correlationID, "play", nn)
	var err error
	log.V(1).Info("update status")

	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		p.Status.State = state
		p.Status.Message = message

		meta.SetStatusCondition(&p.Status.Conditions, metav1.Condition{
			Type:    string(state),
			Status:  r.GetStatus(state),
			Reason:  string(state),
			Message: message,
		})
		if err := r.Status().Update(ctx, p); err != nil {
			log.Error(err, "unable to update play status (retry)")
			return err
		}
		return nil
	})
	if err != nil {
		log.Error(err, "unable to update play status")
		return err
	}
	return nil
}

func (r *PlayReconciler) GetStatus(state ci.State) metav1.ConditionStatus {
	switch state {
	case ci.Errored, ci.Cancelled, ci.Failed:
		return metav1.ConditionFalse
	case ci.Succeeded:
		return metav1.ConditionTrue
	default:
		return metav1.ConditionUnknown
	}
}

// UpdateName set the status of tekton resource state
func (r *PlayReconciler) UpdateName(ctx context.Context, p *ci.Play, name string) error {
	correlationID, _ := util.GetCorrelationIDFromContext(ctx)
	nn, _ := util.GetPlayFromContext(ctx)
	log := ctrl.Log.WithName("Reconcile").WithName("UpdateName").WithValues("correlation_id", correlationID, "play", nn)
	var err error
	log.V(1).Info("update name")
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		p.Status.PipelineRunName = name
		if err := r.Status().Update(ctx, p); err != nil {
			log.Error(err, "unable to update play name (retry)")
			return err
		}
		return nil
	})
	if err != nil {
		log.Error(err, "unable to update play name")
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ci.Play{}).
		Owns(&tkn.PipelineRun{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}

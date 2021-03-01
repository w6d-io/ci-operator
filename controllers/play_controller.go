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

package controllers

import (
	"context"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/tekton/pipelinerun"
	"github.com/w6d-io/ci-operator/internal/util"
	"github.com/w6d-io/ci-operator/pkg/play"
	"github.com/w6d-io/ci-operator/pkg/webhook"
	"k8s.io/apimachinery/pkg/runtime"
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

// +kubebuilder:rbac:groups=ci.w6d.io,resources=plays,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ci.w6d.io,resources=plays/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tekton.dev,resources=pipelineresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tekton.dev,resources=pipelineresources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tekton.dev,resources=pipelineruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tekton.dev,resources=pipelineruns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tekton.dev,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tekton.dev,resources=pipelines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tekton.dev,resources=taskruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tekton.dev,resources=taskruns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tekton.dev,resources=tasks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tekton.dev,resources=tasks/status,verbs=get;update;patch

func (r *PlayReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithName("Reconcile").WithValues("play", req.NamespacedName)
	// get the play resource
	p := new(ci.Play)

	if err := r.Get(ctx, req.NamespacedName, p); err != nil {
		log.Error(err, "unable to fetch Play")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	p.Spec.Name = strings.ToLower(p.Spec.Name)
	p.Spec.Environment = strings.ToLower(p.Spec.Environment)
	log = log.WithValues("cx-namespace", util.InNamespace(p))
	log.V(1).Info("req name " + req.Name)
	log.V(1).Info("get pipelinerun " + util.GetCINamespacedName(pipelinerun.Prefix, p).String())
	var childPr tkn.PipelineRun
	err := r.Get(ctx, util.GetCINamespacedName(pipelinerun.Prefix, p), &childPr)
	if client.IgnoreNotFound(err) != nil {
		log.Error(err, "Unable to list child PipelineRuns")
		return ctrl.Result{}, err
	}

	if !apierrors.IsNotFound(err) {
		if childPr.Name != "" && p.Status.PipelineRunName != childPr.Name {
			p.Status.PipelineRunName = childPr.Name
			log.V(1).Info("updating play status")
			if err := r.Status().Update(ctx, p); err != nil {
				log.Error(err, "unable to update Play status")
				return ctrl.Result{}, err
			}
		}
		if util.Condition(childPr.Status.Conditions) != p.Status.State {
			p.Status.State = util.Condition(childPr.Status.Conditions)
			log.V(2).Info("update status", "status", p.Status.State,
				"step", "1")
			p.Status.State = ci.Succeeded
			if err := r.Status().Update(ctx, p); err != nil {
				log.Error(err, "unable to update Play status")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	log.V(1).Info("pipelinerun not found")
	log.V(1).Info("getting all pipeline run")
	var prs tkn.PipelineRunList
	if err := r.List(ctx, &prs, util.InNamespace(p)); util.IgnoreNotExists(err) != nil {
		log.Error(err, "Unable to list PipelineRuns in ", util.InNamespace(p))
		p.Status.State = ci.Errored
		log.V(2).Info("update status", "status", p.Status.State,
			"step", "2")
		if err := r.Status().Update(ctx, p); err != nil {
			return ctrl.Result{}, err
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
	log.V(1).Info("pipelinerun", "running", len(runningPipeline),
		"cx-namespace", util.InNamespace(p))

	log.V(1).Info("get limitCi")
	var limits ci.LimitCiList
	if err := r.List(ctx, &limits, util.InNamespace(p)); client.IgnoreNotFound(err) != nil {
		log.Error(err, "unable to list LimitCi")
		p.Status.State = ci.Errored
		log.V(2).Info("update status", "status", p.Status.State,
			"step", "3")
		if err := r.Status().Update(ctx, p); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	if len(limits.Items) > 0 && (limits.Items[0].Spec.Concurrent <= int64(len(runningPipeline))) {
		log.V(1).Info("limit ci", "action", "queued")
		p.Status.State = ci.Queued
		log.V(2).Info("update status", "status", p.Status.State,
			"step", "4")
		if err := r.Status().Update(ctx, p); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 5*time.Second}, nil
		//return ctrl.Result{}, nil
	}
	err = play.CreateCI(ctx, p, r.Log, r, r.Scheme)
	if err != nil {
		log.Error(err, "Failed to create CI")
		p.Status.State = ci.Errored
		log.V(2).Info("update status", "status", p.Status.State,
			"step", "5")
		if err := r.Status().Update(ctx, p); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}
	p.Status.State = ci.Succeeded
	log.V(2).Info("update status", "status", p.Status.State,
		"step", "6")
	if err := r.Status().Update(ctx, p); err != nil {
		return ctrl.Result{}, err
	}

	if err := webhook.BuildPlayPayload(p, ci.Unknown, log); err != nil {
		log.Error(err, "build payload of play")
	}

	payload := webhook.GetPayLoad()
	payload.SetStatus(ci.Succeeded)
	nn := util.GetCINamespacedName("pipeline-run", p)
	payload.SetObjectNamespacedName(nn)
	if err := payload.DoSend(config.GetWebhooks()); err != nil {
		log.Error(err, "webhook")
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

func (r *PlayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ci.Play{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}

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

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	resourcev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/resource/v1alpha1"

	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/w6d-io/ci-operator/controllers"
	"github.com/w6d-io/ci-operator/util"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	civ1alpha1 "github.com/w6d-io/ci-operator/api/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// +kubebuilder:scaffold:imports
)

// Version microservice version
var Version = ""

// Revision git commit
var Revision = ""

// GoVersion ...
var GoVersion = ""

// Built Date built
var Built = ""

// OsArch ...
var OsArch = ""

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {

	_ = clientgoscheme.AddToScheme(scheme)

	_ = civ1alpha1.AddToScheme(scheme)
	_ = tkn.AddToScheme(scheme)
	_ = resourcev1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool

	setupLog.Info("managed flag")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Encoder: zapcore.NewConsoleEncoder(util.TextEncoderConfig()),
	}
	util.BindFlags(&opts, flag.CommandLine)
	flag.Parse()
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	if !seen["config"] {
		fmt.Print("config file is missing\n")
		setupLog.Error(errors.New("flag error"), "config file is missing")
		os.Exit(1)
	}
	if err := controllers.Cfg.Validate(); err != nil {
		fmt.Printf("error : %s\n", err)
		setupLog.Error(err, "config loading error")
		os.Exit(1)
	}
	setupLog.Info("set opts")
	opts.Development = true
	opts.StacktraceLevel = zapcore.PanicLevel
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "2f8df8b9.w6d.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.PlayReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Play"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Play")
		os.Exit(1)
	}
	if err = (&civ1alpha1.Play{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Play")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager", "Version", Version, "Built", Built, "Revision", Revision)
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

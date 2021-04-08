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

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	resourcev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/resource/v1alpha1"
	civ1alpha1 "github.com/w6d-io/ci-operator/api/v1alpha1"
	zapraw "go.uber.org/zap"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/w6d-io/ci-operator/controllers"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

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
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(civ1alpha1.AddToScheme(scheme))
	utilruntime.Must(tkn.AddToScheme(scheme))
	utilruntime.Must(resourcev1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", util.LookupEnvOrString("METRICS_ADDRESS", ":8080"), "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", util.LookupEnvOrString("PROBE_ADDRESS", ":8081"), "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", util.LookupEnvOrBool("ENABLE_LEADER", false),
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development:     os.Getenv("RELEASE") != "prod",
		StacktraceLevel: zapcore.PanicLevel,
		Encoder:         zapcore.NewConsoleEncoder(util.TextEncoderConfig()),
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

	if err := config.Validate(); err != nil {
		fmt.Printf("error : %s\n", err)
		setupLog.Error(err, "config loading error")
		os.Exit(1)
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts), zap.RawZapOpts(zapraw.AddCaller(), zapraw.AddCallerSkip(-1))))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "2f8df8b9.ci.w6d.io",
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

	// +kubebuilder:scaffold:builder
	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err = (&civ1alpha1.Play{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Play")
			os.Exit(1)
		}
	}
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager", "Version", Version, "Built",
		Built, "Revision", Revision, "Arch", OsArch, "GoVersion", GoVersion)
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

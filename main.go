/*


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
	"context"
	"crypto/ecdsa"
	"flag"
	"os"

	"github.com/dlorenc/cosigned/pkg/cosigned"
	"github.com/sigstore/cosign/pkg/cosign/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = corev1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

var secretKeyRef string

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&secretKeyRef, "secret-key-ref", "", "The secret that includes pub/private key pair")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "95b304f4.sigstore.dev",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	hookServer := mgr.GetWebhookServer()
	hookServer.Register("/validate-v1-pod", &webhook.Admission{Handler: &podValidator{Client: mgr.GetClient()}})

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// +kubebuilder:webhook:path=/validate-v1-pod,mutating=false,failurePolicy=ignore,groups="",resources=pods,verbs=create;update,versions=v1,name=cosigned.sigstore.dev

// podValidator validates Pods
type podValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// podValidator admits a pod if a specific annotation exists.
func (v *podValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	if err := v.decoder.Decode(req, pod); err != nil {
		setupLog.Error(err, "decoding", "req", req)
		return admission.Denied("error decoding")
	}

	setupLog.Info("looking for secret", "secret", secretKeyRef)
	cfg, err := kubernetes.GetKeyPairSecret(ctx, secretKeyRef)

	if err != nil {
		return admission.Denied(err.Error())
	}
	if cfg == nil {
		return admission.Denied("no keys configured")
	}

	keys := cosigned.Keys(cfg.Data)
	setupLog.Info("got keys", "cosign.pub", keys)
	for _, c := range pod.Spec.Containers {
		if !valid(ctx, c.Image, keys) {
			return admission.Denied("invalid signatures")
		}
	}
	return admission.Allowed("valid signatures!")
}

func valid(ctx context.Context, img string, keys []*ecdsa.PublicKey) bool {
	for _, k := range keys {
		sps, err := cosigned.Signatures(ctx, img, k)
		if err != nil {
			setupLog.Error(err, "checking signatures", "image", img)
			return false
		}
		if len(sps) > 0 {
			setupLog.Info("valid signatures", "image", img, "key", k)
			return true
		}
	}
	return false
}

// podValidator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (v *podValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

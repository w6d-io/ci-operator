package webhook_test

import (
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/pkg/webhook"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("Webhook", func() {
	var (
		play   *ci.Play
		status ci.State
		logger logr.Logger
	)
	BeforeEach(func() {
		play = &ci.Play{
			TypeMeta: metav1.TypeMeta{
				APIVersion: ci.GroupVersion.Group + "/" + ci.GroupVersion.Version,
			},
		}
		status = ci.Running
		logger = ctrl.Log.WithName("test")
	})
	Describe("Send payload", func() {
		Context("When all resource has been created", func() {
			It("Build payload", func() {
				Expect(webhook.BuildPlayPayload(play, status, logger)).To(BeNil())
			})
			It("Get list of subscribers", func() {})
			It("Send to subscribers", func() {
				Expect(webhook.GetPayLoad().Send("")).To(BeNil())
			})
		})
		Context("When some resource creation failed", func() {
			It("Build payload", func() {})
			It("Get list of subscribers", func() {})
			It("Send to subscribers", func() {})
		})
	})
})

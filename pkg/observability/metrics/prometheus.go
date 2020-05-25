package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	promNamespace        = "kubewebhook"
	promWebhookSubsystem = "admission_webhook"
)

// Prometheus is the implementation of a metrics Recorder for
// Prometheus system.
type Prometheus struct {
	// Metrics.
	admissionReview         *prometheus.CounterVec
	admissionReviewErr      *prometheus.CounterVec
	admissionReviewDuration *prometheus.HistogramVec
	// Validation Metrics
	validationReviewResult *prometheus.CounterVec

	reg prometheus.Registerer
}

// NewPrometheus returns a new Prometheus metrics backend.
func NewPrometheus(registry prometheus.Registerer) *Prometheus {
	p := &Prometheus{
		reg: registry,

		admissionReview: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: promNamespace,
			Subsystem: promWebhookSubsystem,
			Name:      "admission_reviews_total",
			Help:      "Total number of admission reviews handled.",
		}, []string{"webhook", "namespace", "resource", "operation", "kind"}),

		admissionReviewErr: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: promNamespace,
			Subsystem: promWebhookSubsystem,
			Name:      "admission_review_errors_total",
			Help:      "Total number of admission review errors when handling.",
		}, []string{"webhook", "namespace", "resource", "operation", "kind"}),

		admissionReviewDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: promNamespace,
			Subsystem: promWebhookSubsystem,
			Name:      "admission_review_duration_seconds",
			Help:      "The duration of the admission review.",
		}, []string{"webhook", "namespace", "resource", "operation", "kind"}),

		validationReviewResult: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: promNamespace,
			Subsystem: promWebhookSubsystem,
			Name:      "validation_review_results_total",
			Help:      "Total number of validation reviews",
		}, []string{"webhook", "namespace", "resource", "operation", "allowed"}),
	}

	p.registerMetrics()
	return p
}

func (p *Prometheus) registerMetrics() {
	p.reg.MustRegister(
		p.admissionReview,
		p.admissionReviewErr,
		p.admissionReviewDuration,
		p.validationReviewResult,
	)
}

// IncAdmissionReview satisfies Recorder interface.
func (p *Prometheus) IncAdmissionReview(webhook, namespace, resource string, operation Operation, kind ReviewKind) {
	p.admissionReview.WithLabelValues(
		webhook,
		namespace,
		string(resource),
		string(operation),
		string(kind)).Inc()
}

// IncAdmissionReviewError satisfies Recorder interface.
func (p *Prometheus) IncAdmissionReviewError(webhook, namespace, resource string, operation Operation, kind ReviewKind) {
	p.admissionReviewErr.WithLabelValues(
		webhook,
		namespace,
		string(resource),
		string(operation),
		string(kind)).Inc()
}

// ObserveAdmissionReviewDuration satisfies Recorder interface.
func (p *Prometheus) ObserveAdmissionReviewDuration(webhook, namespace, resource string, operation Operation, kind ReviewKind, start time.Time) {
	secs := p.getDuration(start).Seconds()
	p.admissionReviewDuration.WithLabelValues(
		webhook,
		namespace,
		string(resource),
		string(operation),
		string(kind)).Observe(secs)
}

// IncValidationReviewResult satisfies Recorder interface.
func (p *Prometheus) IncValidationReviewResult(webhook, namespace, resource string, operation Operation, allowed bool) {
	p.validationReviewResult.WithLabelValues(
		webhook,
		namespace,
		string(resource),
		string(operation),
		strconv.FormatBool(allowed),
	).Inc()
}

func (p *Prometheus) getDuration(start time.Time) time.Duration {
	return time.Since(start)
}

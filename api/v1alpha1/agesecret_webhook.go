/*
Copyright 2022.

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
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var ageSecretLog = logf.Log.WithName("agesecret-resource")

func (r *AgeSecret) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithDefaulter(&AgeSecretCustomDefaulter{}).
		WithValidator(&AgeSecretCustomValidator{}).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-gitopssecret-snappcloud-io-v1alpha1-agesecret,mutating=true,failurePolicy=fail,sideEffects=None,groups=gitopssecret.snappcloud.io,resources=agesecrets,verbs=create;update,versions=v1alpha1,name=magesecret.kb.io,admissionReviewVersions=v1

// AgeSecretCustomDefaulter mutates AgeSecret resources on create and update.
// +kubebuilder:object:generate=false
type AgeSecretCustomDefaulter struct{}

var _ webhook.CustomDefaulter = &AgeSecretCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (d *AgeSecretCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	ageSecret, ok := obj.(*AgeSecret)
	if !ok {
		return fmt.Errorf("expected an AgeSecret object but got %T", obj)
	}
	ageSecretLog.Info("default", "name", ageSecret.Name)
	return nil
}

//+kubebuilder:webhook:path=/validate-gitopssecret-snappcloud-io-v1alpha1-agesecret,mutating=false,failurePolicy=fail,sideEffects=None,groups=gitopssecret.snappcloud.io,resources=agesecrets,verbs=create;update,versions=v1alpha1,name=vagesecret.kb.io,admissionReviewVersions=v1

// AgeSecretCustomValidator validates AgeSecret resources.
// +kubebuilder:object:generate=false
type AgeSecretCustomValidator struct{}

var _ webhook.CustomValidator = &AgeSecretCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (v *AgeSecretCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	ageSecret, ok := obj.(*AgeSecret)
	if !ok {
		return nil, fmt.Errorf("expected an AgeSecret object but got %T", obj)
	}
	ageSecretLog.Info("validate create", "name", ageSecret.Name)
	return nil, ageSecret.ValidateAgeSecret()
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (v *AgeSecretCustomValidator) ValidateUpdate(_ context.Context, _, newObj runtime.Object) (admission.Warnings, error) {
	ageSecret, ok := newObj.(*AgeSecret)
	if !ok {
		return nil, fmt.Errorf("expected an AgeSecret object but got %T", newObj)
	}
	ageSecretLog.Info("validate update", "name", ageSecret.Name)
	return nil, ageSecret.ValidateAgeSecret()
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (v *AgeSecretCustomValidator) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	ageSecret, ok := obj.(*AgeSecret)
	if !ok {
		return nil, fmt.Errorf("expected an AgeSecret object but got %T", obj)
	}
	ageSecretLog.Info("validate delete", "name", ageSecret.Name)
	return nil, nil
}

func (r *AgeSecret) ValidateAgeSecret() error {
	if strings.TrimSpace(r.Spec.StringData) == "" {
		return fmt.Errorf("stringData can not be empty")
	}

	if strings.TrimSpace(r.Spec.AgeKeyRef) == "" {
		return fmt.Errorf("AgeKey reference can not be empty")
	}
	return nil
}

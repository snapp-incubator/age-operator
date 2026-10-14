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

	"filippo.io/age"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var (
	ageKeyLog    = logf.Log.WithName("agekey-resource")
	ageKeyPrefix = "AGE-SECRET-KEY-"
)

func (r *AgeKey) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, r).
		WithDefaulter(&AgeKeyCustomDefaulter{}).
		WithValidator(&AgeKeyCustomValidator{}).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-gitopssecret-snappcloud-io-v1alpha1-agekey,mutating=true,failurePolicy=fail,sideEffects=None,groups=gitopssecret.snappcloud.io,resources=agekeys,verbs=create;update,versions=v1alpha1,name=magekey.kb.io,admissionReviewVersions=v1

// AgeKeyCustomDefaulter mutates AgeKey resources on create and update.
// +kubebuilder:object:generate=false
type AgeKeyCustomDefaulter struct{}

var _ admission.Defaulter[*AgeKey] = &AgeKeyCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the type.
func (d *AgeKeyCustomDefaulter) Default(_ context.Context, ageKey *AgeKey) error {
	ageKeyLog.Info("default", "name", ageKey.Name)
	return nil
}

//+kubebuilder:webhook:path=/validate-gitopssecret-snappcloud-io-v1alpha1-agekey,mutating=false,failurePolicy=fail,sideEffects=None,groups=gitopssecret.snappcloud.io,resources=agekeys,verbs=create;update,versions=v1alpha1,name=vagekey.kb.io,admissionReviewVersions=v1

// AgeKeyCustomValidator validates AgeKey resources.
// +kubebuilder:object:generate=false
type AgeKeyCustomValidator struct{}

var _ admission.Validator[*AgeKey] = &AgeKeyCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type.
func (v *AgeKeyCustomValidator) ValidateCreate(_ context.Context, ageKey *AgeKey) (admission.Warnings, error) {
	ageKeyLog.Info("validate create", "name", ageKey.Name)
	return nil, ageKey.ValidateAgeKey()
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type.
func (v *AgeKeyCustomValidator) ValidateUpdate(_ context.Context, _, ageKey *AgeKey) (admission.Warnings, error) {
	ageKeyLog.Info("validate update", "name", ageKey.Name)
	return nil, ageKey.ValidateAgeKey()
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type.
func (v *AgeKeyCustomValidator) ValidateDelete(_ context.Context, ageKey *AgeKey) (admission.Warnings, error) {
	ageKeyLog.Info("validate delete", "name", ageKey.Name)
	return nil, nil
}

func (r *AgeKey) ValidateAgeKey() error {
	if !strings.HasPrefix(r.Spec.AgeSecretKey, ageKeyPrefix) {
		return fmt.Errorf("AgeKey must start with %v prefix", ageKeyPrefix)
	}
	if _, err := age.ParseX25519Identity(r.Spec.AgeSecretKey); err != nil {
		ageKeyLog.Info("validate AgeKey", "name", r.Name, "error", err)
		return fmt.Errorf("provided AgeKey is not valid")
	}
	return nil
}

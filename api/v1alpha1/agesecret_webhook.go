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
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

// log is for logging in this package.
var ageSecretLog = logf.Log.WithName("agesecret-resource")

func (r *AgeSecret) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-gitopssecret-snappcloud-io-v1alpha1-agesecret,mutating=true,failurePolicy=fail,sideEffects=None,groups=gitopssecret.snappcloud.io,resources=agesecrets,verbs=create;update,versions=v1alpha1,name=magesecret.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &AgeSecret{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *AgeSecret) Default() {
	ageSecretLog.Info("default", "name", r.Name)

}

//+kubebuilder:webhook:path=/validate-gitopssecret-snappcloud-io-v1alpha1-agesecret,mutating=false,failurePolicy=fail,sideEffects=None,groups=gitopssecret.snappcloud.io,resources=agesecrets,verbs=create;update,versions=v1alpha1,name=vagesecret.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &AgeSecret{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *AgeSecret) ValidateCreate() (admission.Warnings, error) {
	ageSecretLog.Info("validate create", "name", r.Name)
	if err := r.ValidateAgeSecret(); err != nil {
		return nil, err
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AgeSecret) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	ageSecretLog.Info("validate update", "name", r.Name)
	if err := r.ValidateAgeSecret(); err != nil {
		return nil, err
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AgeSecret) ValidateDelete() (admission.Warnings, error) {
	ageSecretLog.Info("validate delete", "name", r.Name)

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

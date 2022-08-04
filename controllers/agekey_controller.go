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

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/snapp-incubator/age-operator/k8sutils"
	"k8s.io/apimachinery/pkg/api/errors"

	gitopssecretsnappcloudiov1alpha1 "github.com/snapp-incubator/age-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	err error
)

// AgeKeyReconciler reconciles a AgeKey object
type AgeKeyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Logger logr.Logger
}

//+kubebuilder:rbac:groups=gitopssecret.snappcloud.io,resources=agekeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gitopssecret.snappcloud.io,resources=agekeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gitopssecret.snappcloud.io,resources=agekeys/finalizers,verbs=update

func (r *AgeKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Logger.WithValues("Request.NamespacedName", req.NamespacedName)
	reqLogger.Info("Reconcile Started")

	ageKeyInstance := &gitopssecretsnappcloudiov1alpha1.AgeKey{}

	err = r.Client.Get(ctx, req.NamespacedName, ageKeyInstance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if err = k8sutils.HandleAgeKeyFinalizers(ageKeyInstance, r.Client); err != nil {
		return ctrl.Result{}, err
	}

	if err = k8sutils.AddAgeKeyFinalizers(ageKeyInstance, r.Client); err != nil {
		return ctrl.Result{}, err
	}

	if err = k8sutils.ValidateAgeKey(ageKeyInstance, r.Client); err != nil {
		return ctrl.Result{}, err
	}

	if err = k8sutils.CreateAgeKeyFile(ageKeyInstance); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AgeKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gitopssecretsnappcloudiov1alpha1.AgeKey{}).
		Complete(r)
}

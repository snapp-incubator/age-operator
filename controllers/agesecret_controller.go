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
	"github.com/snapp-incubator/age-operator/lang"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"time"

	gitopssecretsnappcloudiov1alpha1 "github.com/snapp-incubator/age-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	requeueAfterTime = 1 * time.Minute
)

// AgeSecretReconciler reconciles a AgeSecret object
type AgeSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Logger logr.Logger
}

//+kubebuilder:rbac:groups=gitopssecret.snappcloud.io,resources=agesecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gitopssecret.snappcloud.io,resources=agesecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gitopssecret.snappcloud.io,resources=agesecrets/finalizers,verbs=update

func (r *AgeSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Logger.WithValues("Request.NamespacedName", req.NamespacedName)
	reqLogger.Info("Reconcile Started")

	ageSecretInstance := &gitopssecretsnappcloudiov1alpha1.AgeSecret{}

	if err = r.Client.Get(ctx, req.NamespacedName, ageSecretInstance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	goOn, err := k8sutils.HandleAgeSecretFinalizers(ageSecretInstance, r.Client)
	if err != nil {
		reqLogger.Info("could not handle AgeSecret finalizers", "namespace", ageSecretInstance.GetNamespace(), "name", ageSecretInstance.GetName(), "error", err)
		return ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}
	if !goOn {
		return ctrl.Result{}, nil
	}

	if err = k8sutils.AddAgeSecretFinalizers(ageSecretInstance, r.Client); err != nil {
		reqLogger.Info("could not add required finalizers to AgeSecret", "namespace", ageSecretInstance.GetNamespace(), "name", ageSecretInstance.GetName(), "error", err)
		return ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, err
	}

	if ageSecretInstance.Spec.Suspend {
		reqLogger.Info("AgeSecret suspended", "namespace", ageSecretInstance.GetNamespace(), "name", ageSecretInstance.GetName())
		return ctrl.Result{}, nil
	}

	refAgeKey, errAgeKey := k8sutils.CheckAgeKeyReference(ageSecretInstance, r.Client)
	if errAgeKey != nil {
		reqLogger.Info("could not fetch AgeSecret referenced AgeKey", "namespace", ageSecretInstance.GetNamespace(), "name", ageSecretInstance.GetName(), "error", errAgeKey)
		return ctrl.Result{Requeue: true, RequeueAfter: requeueAfterTime}, errAgeKey
	}

	decryptedStringData, errDecrypt := k8sutils.DecryptAgeSecret(ageSecretInstance, r.Client, refAgeKey)
	if errDecrypt != nil {
		reqLogger.Info("could not decrypt AgeSecret", "namespace", ageSecretInstance.GetNamespace(), "name", ageSecretInstance.GetName(), "error", errDecrypt)
		return ctrl.Result{}, errDecrypt
	}

	err = k8sutils.CreateChildFromAgeSecret(ageSecretInstance, r.Client, decryptedStringData)
	if err != nil {
		reqLogger.Info("could not create child secret from AgeSecret", "namespace", ageSecretInstance.GetNamespace(), "name", ageSecretInstance.GetName(), "error", err)
		return ctrl.Result{}, err
	}

	ageSecretInstance.Status.Health = lang.AgeSecretStatusHealthy
	_ = r.Status().Update(ctx, ageSecretInstance)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AgeSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gitopssecretsnappcloudiov1alpha1.AgeSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

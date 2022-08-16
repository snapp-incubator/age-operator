package k8sutils

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/snapp-incubator/age-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	AgeKeyFinalizer    string = "AgeKeyFinalizer"
	AgeSecretFinalizer string = "AgeSecretFinalizer"
)

var logObj = logf.Log.WithName("controller")

func finalizerLogger(namespace string, name string) logr.Logger {
	reqLogger := logObj.WithValues("Request.Service.Namespace", namespace, "Request.Finalizer.Name", name)
	return reqLogger
}

func NewLogger(namespace string, name string) logr.Logger {
	reqLogger := logObj.WithValues("Request.Namespace", namespace, "Request.Name", name)
	return reqLogger
}

func HandleAgeKeyFinalizers(ageKey *v1alpha1.AgeKey, k8sclient client.Client) error {
	if ageKey.GetDeletionTimestamp() != nil {
		logger := finalizerLogger(ageKey.GetNamespace(), AgeKeyFinalizer)
		if controllerutil.ContainsFinalizer(ageKey, AgeKeyFinalizer) {
			if err := finalizeAgeKey(ageKey); err != nil {
				logger.Error(err, "Could not execute finalizer")
				return err
			}
			controllerutil.RemoveFinalizer(ageKey, AgeKeyFinalizer)
			if err := k8sclient.Update(context.Background(), ageKey); err != nil {
				logger.Error(err, "Could not remove finalizer "+AgeKeyFinalizer)
				return err
			}
		}
	}
	return nil
}

func HandleAgeSecretFinalizers(ageSecret *v1alpha1.AgeSecret, k8sclient client.Client) error {
	if ageSecret.GetDeletionTimestamp() != nil {
		logger := finalizerLogger(ageSecret.GetNamespace(), AgeSecretFinalizer)
		if controllerutil.ContainsFinalizer(ageSecret, AgeSecretFinalizer) {
			if err := finalizeAgeSecret(ageSecret, k8sclient); err != nil {
				return err
			}
			controllerutil.RemoveFinalizer(ageSecret, AgeSecretFinalizer)
			if err := k8sclient.Update(context.Background(), ageSecret); err != nil {
				logger.Error(err, "Could not remove finalizer "+AgeSecretFinalizer)
				return err
			}
		}
	}
	return nil
}

func AddAgeKeyFinalizers(ageKey *v1alpha1.AgeKey, k8sclient client.Client) error {
	if !controllerutil.ContainsFinalizer(ageKey, AgeKeyFinalizer) {
		controllerutil.AddFinalizer(ageKey, AgeKeyFinalizer)
		return k8sclient.Update(context.Background(), ageKey)
	}
	return nil
}

func AddAgeSecretFinalizers(ageSecret *v1alpha1.AgeSecret, k8sclient client.Client) error {
	if !controllerutil.ContainsFinalizer(ageSecret, AgeSecretFinalizer) {
		controllerutil.AddFinalizer(ageSecret, AgeSecretFinalizer)
		return k8sclient.Update(context.Background(), ageSecret)
	}
	return nil
}

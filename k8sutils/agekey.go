package k8sutils

import (
	"context"
	"filippo.io/age"
	"fmt"
	"github.com/snapp-incubator/age-operator/api/v1alpha1"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	AgeKeysRootPath = "/tmp/keys/"
)

func finalizeAgeKey(ageKey *v1alpha1.AgeKey) error {
	logger := finalizerLogger(ageKey.GetNamespace(), AgeKeyFinalizer)
	fullPath := GenerateAgeKeyFullPath(ageKey)
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		logger.Error(err, "Could not get stats of file", ageKey.GetName())
		return fmt.Errorf("could not finalize")
	}
	if err := os.Remove(fullPath); err != nil {
		logger.Error(err, "Could not remove file", ageKey.GetName())
		return fmt.Errorf("could not finalize")
	}
	return nil
}

func CreateAgeKeyFile(ageKey *v1alpha1.AgeKey) error {
	logger := NewLogger(ageKey.GetNamespace(), ageKey.GetName())
	fullPath := GenerateAgeKeyFullPath(ageKey)
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			parentDir := GenerateAgeKeyParentDir(ageKey)
			if errDirCreation := os.MkdirAll(parentDir, os.ModePerm); errDirCreation != nil {
				logger.Error(err, "could not create directory", "path", ageKey)
				return errDirCreation
			}
			if _, errCreate := os.Create(fullPath); errCreate != nil {
				return errCreate
			}
		} else {
			logger.Error(err, "could not stat file for creating")
			return err
		}
	}
	return writeAgeKeyToFile(fullPath, ageKey)
}

func writeAgeKeyToFile(fullPath string, ageKey *v1alpha1.AgeKey) error {
	var file, err = os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.WriteString(strings.TrimSpace(ageKey.Spec.AgeSecretKey))
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}

func ValidateAgeKey(ageKey *v1alpha1.AgeKey, k8sclient client.Client) error {
	logger := NewLogger(ageKey.GetNamespace(), ageKey.GetName())
	if _, err := age.ParseX25519Identity(ageKey.Spec.AgeSecretKey); err != nil {
		logger.Error(err, "invalid agekey on reconcile")
		ageKey.Status.Message = "Invalid AgeKey"
		_ = k8sclient.Status().Update(context.Background(), ageKey)
		return fmt.Errorf("invalid AgeKey on field .Spec.ageSecretKey")
	}
	return nil
}

func GenerateAgeKeyFullPath(ageKey *v1alpha1.AgeKey) string {
	return filepath.Join(GenerateAgeKeyParentDir(ageKey), ageKey.GetName())
}

func GenerateAgeKeyParentDir(ageKey *v1alpha1.AgeKey) string {
	return filepath.Join(AgeKeysRootPath, ageKey.GetNamespace())
}

func UpdateAgeKeyStatus(ageKey *v1alpha1.AgeKey, k8sclient client.Client, msg string) error {
	ageKey.Status.Message = msg
	return k8sclient.Status().Update(context.Background(), ageKey)
}

package k8sutils

import (
	"bytes"
	"context"
	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/snapp-incubator/age-operator/api/v1alpha1"
	"github.com/snapp-incubator/age-operator/lang"
	"io"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

var (
	unwantedAnnotations = []string{
		"kubectl.kubernetes.io/last-applied-configuration",
	}
)

func finalizeAgeSecret(ageSecret *v1alpha1.AgeSecret, k8sclient client.Client) error {
	logger := finalizerLogger(ageSecret.GetNamespace(), AgeSecretFinalizer)
	childSecret := &corev1.Secret{}
	err := k8sclient.Get(context.Background(), types.NamespacedName{Name: ageSecret.GetName(), Namespace: ageSecret.GetNamespace()}, childSecret)
	if err != nil {
		if !errors.IsNotFound(err) {
			logger.Error(err, "Could not get child secret"+ageSecret.GetName())
			return err
		} else {
			return nil
		}
	}
	if err = k8sclient.Delete(context.Background(), childSecret); err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "Could not delete child secret "+ageSecret.GetName())
		return err
	}
	return nil
}

func CreateChildFromAgeSecret(ageSecret *v1alpha1.AgeSecret, k8sclient client.Client, rawStringData map[string]string) error {
	childSecretLabels := cloneLabels(ageSecret.GetLabels())
	childSecretAnnotations := cloneAnnotations(ageSecret.GetAnnotations())
	childSecretMetaObj := metav1.ObjectMeta{
		Name:        ageSecret.GetName(),
		Namespace:   ageSecret.GetNamespace(),
		Labels:      childSecretLabels,
		Annotations: childSecretAnnotations,
	}
	secretObj := &corev1.Secret{
		ObjectMeta: childSecretMetaObj,
		Type:       corev1.SecretTypeOpaque,
		StringData: rawStringData,
	}
	return CreateOrUpdateSecretObj(ageSecret, secretObj, k8sclient)
}

func CreateOrUpdateSecretObj(ageSecret *v1alpha1.AgeSecret, secret *corev1.Secret, k8sclient client.Client) error {
	secretToLoad := &corev1.Secret{}
	logger := NewLogger(ageSecret.GetNamespace(), ageSecret.GetName())
	err := k8sclient.Get(context.Background(), types.NamespacedName{Namespace: secret.Namespace, Name: secret.Name}, secretToLoad)
	if err != nil {
		if errors.IsNotFound(err) {
			errCreateChildSecret := k8sclient.Create(context.Background(), secret)
			if errCreateChildSecret != nil {
				ageSecret.Status.Health = lang.AgeSecretStatusUnhealthy
				ageSecret.Status.Message = "could not create child secret"
				_ = k8sclient.Status().Update(context.Background(), ageSecret)
				return errCreateChildSecret
			}
			ageSecret.Status.Health = lang.AgeSecretStatusHealthy
			_ = k8sclient.Status().Update(context.Background(), ageSecret)
			return nil
		}
		ageSecret.Status.Health = lang.AgeSecretStatusUnhealthy
		ageSecret.Status.Message = "could not fetch child secret"
		_ = k8sclient.Status().Update(context.Background(), ageSecret)
		return err
	}

	if !apiequality.Semantic.DeepEqual(secretToLoad, secret) {
		logger.Info("child secret exists but needs to get refreshed")
		err = k8sclient.Update(context.Background(), secret)
		if err != nil {
			logger.Error(err, "could not refresh child secret")
			ageSecret.Status.Health = lang.AgeSecretStatusUnhealthy
			ageSecret.Status.Message = "could not refresh child secret"
			_ = k8sclient.Status().Update(context.Background(), ageSecret)
			return err
		}
	}
	if ageSecret.Status.Health != lang.AgeSecretStatusHealthy {
		ageSecret.Status.Health = lang.AgeSecretStatusHealthy
		errUpdateHealth := k8sclient.Status().Update(context.Background(), ageSecret)
		logger.Error(errUpdateHealth, "Could not update status of ageSecret", "namespace", ageSecret.GetNamespace(), "name", ageSecret.GetName())
	}
	return nil
}

func DecryptAgeSecret(ageSecret *v1alpha1.AgeSecret, k8sclient client.Client, ageKey *v1alpha1.AgeKey) (map[string]string, error) {
	logger := NewLogger(ageSecret.GetNamespace(), ageSecret.GetName())
	identity, err := age.ParseX25519Identity(ageKey.Spec.AgeSecretKey)
	if err != nil {
		return nil, err
	}
	encryptedData := "-----BEGIN AGE ENCRYPTED FILE-----\n" + strings.TrimSpace(ageSecret.Spec.StringData) + "\n-----END AGE ENCRYPTED FILE-----"
	reader, errDecrypt := age.Decrypt(armor.NewReader(strings.NewReader(encryptedData)), identity)
	if errDecrypt != nil {
		logger.Error(errDecrypt, "could not decrypt with given AgeKey")
		return nil, errDecrypt
	}
	output := &bytes.Buffer{}
	if _, err := io.Copy(output, reader); err != nil {
		return nil, err
	}
	var rawStringData map[string]string
	outputBytes := output.Bytes()
	err = yaml.Unmarshal(outputBytes, &rawStringData)
	if err != nil {
		logger.Error(err, "given yaml structure is not valid")
		ageSecret.Status.Message = "given yaml structure is not valid"
		_ = k8sclient.Update(context.Background(), ageSecret)
		return nil, err
	}
	return rawStringData, nil
}

func CheckAgeKeyReference(ageSecret *v1alpha1.AgeSecret, k8sclient client.Client) (*v1alpha1.AgeKey, error) {
	ageKeyObj := &v1alpha1.AgeKey{}
	err := k8sclient.Get(context.Background(), types.NamespacedName{Namespace: ageSecret.Namespace, Name: ageSecret.Spec.AgeKeyRef}, ageKeyObj)
	if err != nil {
		return nil, err
	}
	return ageKeyObj, nil
}

func cloneLabels(labels map[string]string) map[string]string {
	return cloneMap(labels)
}

func cloneAnnotations(annotations map[string]string) map[string]string {
	tmpAnnotations := cloneMap(annotations)
	for _, annotation := range unwantedAnnotations {
		delete(tmpAnnotations, annotation)
	}
	return tmpAnnotations
}

func cloneMap(oldMap map[string]string) map[string]string {
	newMap := make(map[string]string)

	for key, value := range oldMap {
		newMap[key] = value
	}

	return newMap
}

func UpdateAgeSecretStatus(ageSecret *v1alpha1.AgeSecret, k8sclient client.Client, health, msg string) error {
	ageSecret.Status.Health = health
	ageSecret.Status.Message = msg
	return k8sclient.Status().Update(context.Background(), ageSecret)
}

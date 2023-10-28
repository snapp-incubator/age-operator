package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snapp-incubator/age-operator/api/v1alpha1"
	"github.com/snapp-incubator/age-operator/consts"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	"path/filepath"
	"time"
)

var (
	fooValidAgeKeyPath3   = filepath.Join("..", "config", "samples", "_v1alpha1_agekey3.yaml")
	fooValidAgeKeyPath2   = filepath.Join("..", "config", "samples", "_v1alpha1_agekey2.yaml")
	fooValidAgeSecretPath = filepath.Join("..", "config", "samples", "_v1alpha1_agesecret.yaml")
)

var _ = Describe("", func() {
	ctx := context.Background()
	validAgeKeyObj := &v1alpha1.AgeKey{}
	invalidAgeKeyObj := &v1alpha1.AgeKey{}
	validAgeSecretObj := &v1alpha1.AgeSecret{}
	namespaceObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-age-secret",
		},
	}

	BeforeEach(func() {
		content, err := os.ReadFile(fooValidAgeKeyPath2)
		Expect(err).Should(BeNil())

		obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(content, nil, nil)
		validAgeKeyObj = obj.(*v1alpha1.AgeKey)
		Expect(err).Should(BeNil())
	})

	BeforeEach(func() {
		content, err := os.ReadFile(fooValidAgeKeyPath3)
		Expect(err).Should(BeNil())

		obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(content, nil, nil)
		invalidAgeKeyObj = obj.(*v1alpha1.AgeKey)
		Expect(err).Should(BeNil())
	})

	BeforeEach(func() {
		content, err := os.ReadFile(fooValidAgeSecretPath)
		Expect(err).Should(BeNil())

		obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(content, nil, nil)
		validAgeSecretObj = obj.(*v1alpha1.AgeSecret)
		Expect(err).Should(BeNil())
	})

	Context("When creating AgeSecret", func() {
		It("should pass if AgeSecret and referenced AgeKey is valid", func() {
			err = k8sClient.Create(ctx, namespaceObj)
			Expect(err).To(BeNil())

			err = k8sClient.Create(ctx, validAgeKeyObj)
			Expect(err).To(BeNil())

			err = k8sClient.Create(ctx, validAgeSecretObj)
			Expect(err).To(BeNil())
			time.Sleep(time.Second * 3)

			fooSecretObj := &corev1.Secret{}
			err = k8sClient.Get(ctx, types.NamespacedName{Namespace: validAgeSecretObj.Namespace, Name: validAgeSecretObj.Name}, fooSecretObj)
			Expect(err).To(BeNil())
			Expect(fooSecretObj.GetAnnotations()).Should(Equal(validAgeSecretObj.GetAnnotations()))

			// make sure unwanted label is removed
			unwantedLabelExists := false
			secretLabels := fooSecretObj.GetLabels()
			for _, label := range secretLabels {
				for _, unwantedLabel := range consts.ExcessLabels {
					if label == unwantedLabel {
						unwantedLabelExists = true
						break
					}
				}
				if unwantedLabelExists {
					break
				}
			}
			Expect(unwantedLabelExists).To(BeFalse())

			sampleKeyValue, exists := fooSecretObj.Data["sample_key"]
			Expect(string(sampleKeyValue)).Should(Equal("sample_value"))
			Expect(exists).To(BeTrue())

			testKeyValue, exists := fooSecretObj.Data["test_key"]
			Expect(string(testKeyValue)).Should(Equal("test_value"))
			Expect(exists).To(BeTrue())

			err = k8sClient.Delete(ctx, validAgeKeyObj)
			Expect(err).To(BeNil())

			err = k8sClient.Delete(ctx, validAgeSecretObj)
			Expect(err).To(BeNil())
			time.Sleep(time.Second * 3)

			fooSecretObj2 := &corev1.Secret{}
			err = k8sClient.Get(ctx, types.NamespacedName{Namespace: validAgeSecretObj.Namespace, Name: validAgeSecretObj.Name}, fooSecretObj2)
			Expect(err).NotTo(BeNil())
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})

		It("should fail if AgeSecret is valid but referenced AgeKey is not the recipient", func() {
			err = k8sClient.Create(ctx, invalidAgeKeyObj)
			Expect(err).To(BeNil())

			err = k8sClient.Create(ctx, validAgeSecretObj)
			Expect(err).To(BeNil())
			time.Sleep(time.Second * 3)

			fooSecretObj := &corev1.Secret{}
			err = k8sClient.Get(ctx, types.NamespacedName{Namespace: validAgeSecretObj.Namespace, Name: validAgeSecretObj.Name}, fooSecretObj)
			Expect(err).NotTo(BeNil())
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})
	})
})

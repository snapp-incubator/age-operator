package v1alpha1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("", func() {
	const (
		fooName          = "foo-agesecret"
		fooNamespace     = "default"
		sampleStringData = "fake data"
		sampleAgeKey     = "sample-agekey"
	)
	var (
		err error
	)
	fooAgeSecretTypeMeta := metav1.TypeMeta{
		APIVersion: "gitopssecret.snappcloud.io/v1alpha1",
		Kind:       "AgeSecret",
	}
	fooAgeSecretObjectMeta := metav1.ObjectMeta{
		Name:      fooName,
		Namespace: fooNamespace,
	}
	fooAgeSecretMeta := &AgeSecret{
		TypeMeta:   fooAgeSecretTypeMeta,
		ObjectMeta: fooAgeSecretObjectMeta,
	}

	AfterEach(func() {
		err = k8sClient.Delete(ctx, fooAgeSecretMeta)
		if err != nil {
			Expect(errors.IsNotFound(err)).Should(BeTrue())
		}
	})

	Context("When creating AgeSecret", func() {
		It("should fail if age_key_ref is empty", func() {
			fooAgeSecret := &AgeSecret{
				TypeMeta:   fooAgeSecretMeta.TypeMeta,
				ObjectMeta: fooAgeSecretMeta.ObjectMeta,
				Spec: AgeSecretSpec{
					AgeKeyRef:  "",
					StringData: sampleStringData,
				},
			}
			err = k8sClient.Create(ctx, fooAgeSecret)
			Expect(err).NotTo(BeNil())
		})

		It("should fail if stringData is empty", func() {
			fooAgeSecret := &AgeSecret{
				TypeMeta:   fooAgeSecretMeta.TypeMeta,
				ObjectMeta: fooAgeSecretMeta.ObjectMeta,
				Spec: AgeSecretSpec{
					AgeKeyRef:  sampleAgeKey,
					StringData: "",
				},
			}
			err = k8sClient.Create(ctx, fooAgeSecret)
			Expect(err).NotTo(BeNil())
		})

		It("should pass if AgeSecret is valid", func() {
			fooAgeSecret := &AgeSecret{
				TypeMeta:   fooAgeSecretMeta.TypeMeta,
				ObjectMeta: fooAgeSecretMeta.ObjectMeta,
				Spec: AgeSecretSpec{
					AgeKeyRef:  sampleAgeKey,
					StringData: sampleStringData,
				},
			}
			err = k8sClient.Create(ctx, fooAgeSecret)
			Expect(err).To(BeNil())

			fooAgeSecret2 := &AgeSecret{}
			err = k8sClient.Get(ctx, types.NamespacedName{Namespace: fooAgeSecret.GetNamespace(), Name: fooAgeSecret.GetName()}, fooAgeSecret2)
			Expect(err).To(BeNil())
			Expect(fooAgeSecret2.Spec.Suspend).To(BeFalse())
		})
	})
})

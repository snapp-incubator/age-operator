package v1alpha1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("", func() {
	const (
		fooName      = "foo-agekey"
		fooNamespace = "default"
	)
	var (
		err error
	)
	fooAgeKeyTypeMeta := metav1.TypeMeta{
		APIVersion: "gitopssecret.snappcloud.io/v1alpha1",
		Kind:       "AgeKey",
	}
	fooAgeKeyObjectMeta := metav1.ObjectMeta{
		Name:      fooName,
		Namespace: fooNamespace,
	}
	fooAgeKeyMeta := &AgeKey{
		TypeMeta:   fooAgeKeyTypeMeta,
		ObjectMeta: fooAgeKeyObjectMeta,
	}

	AfterEach(func() {
		err = k8sClient.Delete(ctx, fooAgeKeyMeta)
		if err != nil {
			Expect(errors.IsNotFound(err)).Should(BeTrue())
		}
	})

	Context("When creating AgeKey", func() {
		It("should fail if AgeKey is empty", func() {
			fooAgeKey := &AgeKey{
				TypeMeta:   fooAgeKeyMeta.TypeMeta,
				ObjectMeta: fooAgeKeyMeta.ObjectMeta,
				Spec: AgeKeySpec{
					AgeSecretKey: "",
				},
			}
			err = k8sClient.Create(ctx, fooAgeKey)
			Expect(err).NotTo(BeNil())
		})

		It("should fail if AgeKey only has required prefix", func() {
			fooAgeKey := &AgeKey{
				TypeMeta:   fooAgeKeyMeta.TypeMeta,
				ObjectMeta: fooAgeKeyMeta.ObjectMeta,
				Spec: AgeKeySpec{
					AgeSecretKey: ageKeyPrefix,
				},
			}
			err = k8sClient.Create(ctx, fooAgeKey)
			Expect(err).NotTo(BeNil())
		})

		It("should pass if AgeKey is valid", func() {
			fooAgeKey := &AgeKey{
				TypeMeta:   fooAgeKeyMeta.TypeMeta,
				ObjectMeta: fooAgeKeyMeta.ObjectMeta,
				Spec: AgeKeySpec{
					AgeSecretKey: "AGE-SECRET-KEY-1MDTY43V4JSAQF8LRVQ6698JKQJ56XH9JTUXJR7GRQ7CPFVAZQ5GQ5LKNLG",
				},
			}
			err = k8sClient.Create(ctx, fooAgeKey)
			Expect(err).To(BeNil())
		})
	})
})

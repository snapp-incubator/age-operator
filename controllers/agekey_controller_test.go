package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snapp-incubator/age-operator/api/v1alpha1"
	"github.com/snapp-incubator/age-operator/k8sutils"
	"io/ioutil"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	"path/filepath"
	"time"
)

var (
	fooValidAgeKeyPath   = filepath.Join("..", "config", "samples", "_v1alpha1_agekey.yaml")
	fooInvalidAgeKeyPath = filepath.Join("..", "config", "samples", "_v1alpha1_agekey_invalid.yaml")
)

var _ = Describe("Testing AgeKey", func() {
	ctx := context.Background()
	validAgeKeyObj := &v1alpha1.AgeKey{}
	invalidAgeKeyObj := &v1alpha1.AgeKey{}

	BeforeEach(func() {
		content, err := ioutil.ReadFile(fooValidAgeKeyPath)
		Expect(err).Should(BeNil())

		obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(content, nil, nil)
		validAgeKeyObj = obj.(*v1alpha1.AgeKey)
		Expect(err).Should(BeNil())
	})

	BeforeEach(func() {
		content, err := ioutil.ReadFile(fooInvalidAgeKeyPath)
		Expect(err).Should(BeNil())

		obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(content, nil, nil)
		invalidAgeKeyObj = obj.(*v1alpha1.AgeKey)
		Expect(err).Should(BeNil())
	})

	Context("When creating AgeKey", func() {
		It("should fail if AgeKey is invalid", func() {
			err = k8sClient.Create(ctx, invalidAgeKeyObj)
			Expect(err).To(BeNil())
			err = k8sClient.Delete(ctx, invalidAgeKeyObj)
			Expect(err).To(BeNil())
		})

		It("should pass if AgeKey is valid", func() {
			err = k8sClient.Create(ctx, validAgeKeyObj)
			Expect(err).To(BeNil())
			time.Sleep(time.Second * 2)

			// check that file is created successfully
			_, err = os.Stat(k8sutils.GenerateAgeKeyFullPath(validAgeKeyObj))
			Expect(err).To(BeNil())

			err = k8sClient.Delete(ctx, validAgeKeyObj)
			Expect(err).To(BeNil())
			time.Sleep(time.Second * 2)

			// check that file is removed successfully
			_, err = os.Stat(k8sutils.GenerateAgeKeyFullPath(validAgeKeyObj))
			Expect(os.IsNotExist(err)).To(BeTrue())
		})
	})
})

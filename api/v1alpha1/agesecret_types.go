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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AgeSecretSpec defines the desired state of AgeSecret
type AgeSecretSpec struct {
	// +kubebuilder:validation:Required
	AgeKeyRef  string `json:"age_key_ref"`
	StringData string `json:"stringData"`
	// +kubebuilder:validation:Optional
	Suspend bool `json:"suspend"`
}

// AgeSecretStatus defines the observed state of AgeSecret
type AgeSecretStatus struct {
	Health  string `json:"health"`
	Message string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Health",type=string,JSONPath=`.status.health`
//+kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`
//+kubebuilder:printcolumn:name="Suspended",type=string,JSONPath=`.spec.suspend`
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// AgeSecret is the Schema for the agesecrets API
type AgeSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgeSecretSpec   `json:"spec,omitempty"`
	Status AgeSecretStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AgeSecretList contains a list of AgeSecret
type AgeSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AgeSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AgeSecret{}, &AgeSecretList{})
}

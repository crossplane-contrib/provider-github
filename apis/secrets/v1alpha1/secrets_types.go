/*
Copyright 2021 The Crossplane Authors.

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

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// SecretsParameters defines the desired state of a GitHub Secrets.
type SecretsParameters struct {
	// The name of the Repository owner.
	Owner string `json:"owner"`
	// The name of the repository.
	Repository string `json:"repository"`
	// The value of the secret
	Value xpv1.SecretKeySelector `json:"value"`
}

// SecretsObservation are the observable fields of a Secrets.
type SecretsObservation struct {
	// The encrypted value stored
	// +optional
	EncryptValue *string `json:"encrypt_value,omitempty"`
	// Last updated time in Repo Secret GitHub
	// +optional
	LastUpdate *string `json:"last_update,omitempty"`
}

// A SecretsSpec defines the desired state of a Secrets.
type SecretsSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       SecretsParameters `json:"forProvider"`
}

// A SecretsStatus represents the observed state of a Secrets.
type SecretsStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          SecretsObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Secrets is a managed resource that represents a GitHub Secrets
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".status.atProvider.htmlUrl"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,github}
type Secrets struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretsSpec   `json:"spec"`
	Status SecretsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecretsList contains a list of Secrets
type SecretsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Secrets `json:"items"`
}

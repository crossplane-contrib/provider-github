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

// ContentParameters defines the desired state of a GitHub Repository Content.
type ContentParameters struct {
	// The name of the Repository owner.
	// The owner can be an organization or an user.
	Owner string `json:"owner"`

	// The name of the Repository.
	// TODO: Use Selector pattern
	Repo string `json:"repo"`

	// The file path.
	Path string `json:"path"`

	// The commit message.
	Message string `json:"message"`

	// The file content.
	Content string `json:"content"`

	// The branch name. Default to the respository's default branch.
	// +optional
	Branch *string `json:"branch,omitempty"`
}

// ContentSpec defines the desired state of a Content.
type ContentSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ContentParameters `json:"forProvider"`
}

// ContentObservation is the representation of the current state that is observed
type ContentObservation struct {
	URL     string `json:"url,omitempty"`
	HTMLURL string `json:"htmlUrl,omitempty"`
	Sha     string `json:"sha,omitempty"`
}

// ContentStatus represents the observed state of a Content.
type ContentStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ContentObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// A Content is a managed resource that represents a GitHub Repository Content Path
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".status.atProvider.htmlUrl"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,github}
type Content struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContentSpec   `json:"spec"`
	Status ContentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ContentList contains a list of Content
type ContentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Content `json:"items"`
}

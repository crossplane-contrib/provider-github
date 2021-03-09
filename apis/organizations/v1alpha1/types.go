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

// MembershipParameters represents the status of a user's membership in an
// organization or team.
type MembershipParameters struct {
	// GitHub user ID for the person you are inviting. Not required if you provide Email.
	// +optional
	InviteeID *int64 `json:"inviteeId,omitempty"`

	// Email address of the person you are inviting, which can be an existing GitHub user.
	// Not required if you provide InviteeID
	// +optional
	Email *string `json:"email,omitempty"`

	// User is the username of the github user.
	User string `json:"user,omitempty"`

	// Specify role for new member. Can be one of:
	// * admin - Organization owners with full administrative rights to the
	// 	 organization and complete access to all repositories and teams.
	// * direct_member - Non-owner organization members with ability to see
	//   other members and join teams by invitation.
	// * billing_manager - Non-owner organization members with ability to
	//   manage the billing settings of your organization.
	// Default is "direct_member".
	// +optional
	Role *string `json:"role,omitempty"`

	// Name of the organization.
	Organization string `json:"organization"`
}

// MembershipSpec defines the desired state of a Membership.
type MembershipSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       MembershipParameters `json:"forProvider"`
}

// MembershipObservation is the representation of the current state that is observed
type MembershipObservation struct {
	URL *string `json:"url,omitempty"`

	// State is the user's status within the organization or team.
	// Possible values are: "active", "pending"
	State *string `json:"state,omitempty"`

	// TODO(hasheddan): User and Organization are omitted here because they are
	// overly verbose.
}

// MembershipStatus represents the observed state of a Membership.
type MembershipStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          MembershipObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// A Membership is a managed resource that represents a AWS Simple Membership
// +kubebuilder:printcolumn:name="ARN",type="string",JSONPath=".status.atProvider.arn"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,github}
type Membership struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MembershipSpec   `json:"spec"`
	Status MembershipStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MembershipList contains a list of Membership
type MembershipList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Membership `json:"items"`
}

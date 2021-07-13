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

// RepositoryParameters defines the desired state of a GitHub Repository.
type RepositoryParameters struct {
	// The name of the Repository owner.
	// The owner can be an organization or an user.
	Owner string `json:"owner"`

	// The name of the repository.
	Name string `json:"name"`

	// The name of the organization that owns the Repository.
	// +optional
	Organization *string `json:"org,omitempty"`

	// A short description of the repository.
	// +optional
	Description *string `json:"description,omitempty"`

	// A URL with more information about the repository.
	// +optional
	Homepage *string `json:"homepage,omitempty"`

	// Whether the repository is private.
	// Must match with Visibility field.
	// Default: false
	// +optional
	Private *bool `json:"private,omitempty"`

	// Can be public or private. You cannot have private and visibility fields
	// contradictory to each other.
	// If your organization is associated with an enterprise account
	// using GitHub Enterprise Cloud or GitHub Enterprise  Server 2.20+,
	// visibility can also be internal.
	//
	// +optional
	// +kubebuilder:validation:Enum=public;private
	Visibility *string `json:"visibility,omitempty"`

	// Either true to enable issues for this repository or false to
	// disable them.
	// Default: true
	// +optional
	HasIssues *bool `json:"hasIssues,omitempty"`

	// Either true to enable projects for this repository or false
	// to disable them.
	// Note: For organizations that has disabled repository projects, the
	// default is false, and if you pass true, the API returns an error.
	// Default: true
	// +optional
	HasProjects *bool `json:"hasProjects,omitempty"`

	// Either true to enable the wiki for this repository or false
	// to disable it.
	// Default: true
	// +optional
	HasWiki *bool `json:"hasWiki,omitempty"`

	// Either true to make this repo available as a template repository
	// or false to prevent it.
	// Default: false
	// +optional
	IsTemplate *bool `json:"isTemplate,omitempty"`

	// The id of the team that will be granted access to this repository.
	// This is only valid when creating a repository in an organization.
	// +optional
	TeamID *int64 `json:"teamId,omitempty"`

	// Pass true to create an initial commit with empty README.
	// +optional
	AutoInit *bool `json:"autoInit,omitempty"`

	// Desired language or platform .gitignore template to apply.
	// Use the name of the template without the extension.
	// Example: "Haskell".
	// +optional
	GitignoreTemplate *string `json:"gitignoreTemplate,omitempty"`

	// Choose an open source license template that best suits your needs,
	// and then use the license keyword as the license template string.
	// Example: "mpl-2.0".
	// +optional
	LicenseTemplate *string `json:"licenseTemplate,omitempty"`

	// Either true to allow squash-merging pull requests, or false to
	// prevent squash-merging.
	// Default: true
	// +optional
	AllowSquashMerge *bool `json:"allowSquashMerge,omitempty"`

	// Either true to allow merging pull requests with a merge commit,
	// or false to prevent merging pull requests with merge commits.
	// Default: true
	// +optional
	AllowMergeCommit *bool `json:"allowMergeCommit,omitempty"`

	// Either true to allow rebase-merging pull requests, or false to
	// prevent rebase-merging.
	// Default: true
	// +optional
	AllowRebaseMerge *bool `json:"allowRebaseMerge,omitempty"`

	// Either true to allow automatically deleting head branches when
	// pull requests are merged, or false to prevent automatic deletion.
	// Default: false
	// +optional
	DeleteBranchOnMerge *bool `json:"deleteBranchOnMerge,omitempty"`

	// Either true to enable pages for this repository or false
	// to disable it.
	// Default: false
	// +optional
	HasPages *bool `json:"hasPages,omitempty"`

	// Either true to enable downloads for this repository or false
	// to disable it.
	// Default: true
	// +optional
	HasDownloads *bool `json:"hasDownloads,omitempty"`

	// Name of the default branch
	// The branch must already exist in the repository.
	// +optional
	DefaultBranch *string `json:"defaultBranch,omitempty"`

	// True to archive this repository.
	// Note: You cannot unarchive repositories through the API.
	// +optional
	Archived *bool `json:"archived,omitempty"`

	// Reference to the repository template that this
	// repository will be derived from.
	// It is in the format <repository-owner>/<repository-name>
	// (e.g crossplane/provider-github)
	// +optional
	// +immutable
	Template *xpv1.Reference `json:"templateRef,omitempty"`
}

// RepositorySpec defines the desired state of a Repository.
type RepositorySpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       RepositoryParameters `json:"forProvider"`
}

// RepositoryObservation is the representation of the current state that is observed
type RepositoryObservation struct {
	// The ID of the Repository
	ID int64 `json:"id,omitempty"`

	// The NodeID of the Repository
	NodeID string `json:"nodeId,omitempty"`

	// The repository fullname
	// The format is {owner}/{repository_name}
	FullName string `json:"fullName,omitempty"`

	// The name of the repository returned by the API.
	// This field is on the Observation struct to enable update in the repository name
	Name string `json:"name"`

	// TODO: Owner

	// Related Repository URLs

	URL              string `json:"url,omitempty"`
	ArchiveURL       string `json:"archiveUrl,omitempty"`
	AssigneesURL     string `json:"assigneesUrl,omitempty"`
	BlobsURL         string `json:"blobsUrl,omitempty"`
	BranchesURL      string `json:"branchesUrl,omitempty"`
	CollaboratorsURL string `json:"collaboratorsUrl,omitempty"`
	CommentsURL      string `json:"commentsUrl,omitempty"`
	CommitsURL       string `json:"commitsUrl,omitempty"`
	CompareURL       string `json:"compareUrl,omitempty"`
	ContentsURL      string `json:"contentsUrl,omitempty"`
	ContributorsURL  string `json:"contributorsUrl,omitempty"`
	DeploymentsURL   string `json:"deploymentsUrl,omitempty"`
	DownloadsURL     string `json:"downloadsUrl,omitempty"`
	EventsURL        string `json:"eventsUrl,omitempty"`
	ForksURL         string `json:"forksUrl,omitempty"`
	GitCommitsURL    string `json:"gitCommitsUrl,omitempty"`
	GitRefsURL       string `json:"gitRefsUrl,omitempty"`
	GitTagsURL       string `json:"gitTagsUrl,omitempty"`
	HooksURL         string `json:"hooksUrl,omitempty"`
	IssueCommentURL  string `json:"issueCommentUrl,omitempty"`
	IssueEventsURL   string `json:"issueEventsUrl,omitempty"`
	IssuesURL        string `json:"issuesUrl,omitempty"`
	KeysURL          string `json:"keysUrl,omitempty"`
	LabelsURL        string `json:"labelsUrl,omitempty"`
	LanguagesURL     string `json:"languagesUrl,omitempty"`
	MergesURL        string `json:"mergesUrl,omitempty"`
	MilestonesURL    string `json:"milestonesUrl,omitempty"`
	NotificationsURL string `json:"notificationsUrl,omitempty"`
	PullsURL         string `json:"pullsUrl,omitempty"`
	ReleasesURL      string `json:"releasesUrl,omitempty"`
	StargazersURL    string `json:"stargazersUrl,omitempty"`
	StatusesURL      string `json:"statusesUrl,omitempty"`
	SubscribersURL   string `json:"subscribersUrl,omitempty"`
	SubscriptionURL  string `json:"subscriptionUrl,omitempty"`
	TagsURL          string `json:"tagsUrl,omitempty"`
	TreesURL         string `json:"treesUrl,omitempty"`
	TeamsURL         string `json:"teamsUrl,omitempty"`
	HTMLURL          string `json:"htmlUrl,omitempty"`
	CloneURL         string `json:"cloneUrl,omitempty"`
	GitURL           string `json:"gitUrl,omitempty"`
	MirrorURL        string `json:"mirrorUrl,omitempty"`
	SSHURL           string `json:"sshUrl,omitempty"`
	SVNURL           string `json:"svnUrl,omitempty"`

	// Related counters of the repository.

	ForksCount       int `json:"forksCount,omitempty"`
	NetworkCount     int `json:"networkCount,omitempty"`
	OpenIssuesCount  int `json:"openIssuesCount,omitempty"`
	StargazersCount  int `json:"stargazersCount,omitempty"`
	SubscribersCount int `json:"subscribersCount,omitempty"`
	WatchersCount    int `json:"watchersCount,omitempty"`

	// Time that the Repository was created.
	CreatedAt *metav1.Time `json:"createdAt,omitempty"`

	// Time that the Repository was pushed.
	PushedAt *metav1.Time `json:"pushedAt,omitempty"`

	// Time that the Repository was updated.
	UpdatedAt *metav1.Time `json:"updatedAt,omitempty"`

	// Main programming language of the Repository.
	Language string `json:"language,omitempty"`

	// Whether the repository is a fork.
	Fork   bool     `json:"fork,omitempty"`
	Size   int      `json:"size,omitempty"`
	Topics []string `json:"topics,omitempty"`

	// Whether the repository is disabled.
	Disabled    bool            `json:"disabled,omitempty"`
	Permissions map[string]bool `json:"permissions,omitempty"`

	// TODOs below are overly verbose
	// TODO: Parent repository
	// TODO: Source repository
	// TODO: Organization
}

// RepositoryStatus represents the observed state of a Repository.
type RepositoryStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          RepositoryObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// A Repository is a managed resource that represents a GitHub Repository
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".status.atProvider.htmlUrl"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,github}
type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RepositorySpec   `json:"spec"`
	Status RepositoryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RepositoryList contains a list of Repository
type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repository `json:"items"`
}

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

package repositories

import (
	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/go-github/v33/github"
	"github.com/mitchellh/copystructure"
	"github.com/pkg/errors"
)

const (
	errCheckUpToDate = "unable to determine if external resource is up to date"
)

// IsUpToDate checks whether Repository is configured with given RepositoryParameters.
func IsUpToDate(rp *v1alpha1.RepositoryParameters, observed *github.Repository) (bool, error) {
	generated, err := copystructure.Copy(observed)
	if err != nil {
		return true, errors.Wrap(err, errCheckUpToDate)
	}
	desired, ok := generated.(*github.Repository)
	if !ok {
		return true, errors.New(errCheckUpToDate)
	}

	GenerateRepository(*rp, desired)

	return cmp.Equal(
		desired,
		observed,
		cmpopts.IgnoreFields(github.Repository{}, "AutoInit"),
	), nil
}

// GenerateRepository produces *github.Repository from RepositoryParameters
func GenerateRepository(rp v1alpha1.RepositoryParameters, r *github.Repository) { // nolint:gocyclo
	if len(rp.Name) != 0 {
		r.Name = ghclient.StringPtr(rp.Name)
	}
	if rp.Description != nil {
		r.Description = rp.Description
	}
	if rp.Homepage != nil {
		r.Homepage = rp.Homepage
	}
	if rp.Private != nil {
		r.Private = rp.Private
	}
	if rp.Visibility != nil {
		r.Visibility = rp.Visibility
	}
	if rp.HasIssues != nil {
		r.HasIssues = rp.HasIssues
	}
	if rp.HasProjects != nil {
		r.HasProjects = rp.HasProjects
	}
	if rp.HasWiki != nil {
		r.HasWiki = rp.HasWiki
	}
	if rp.AutoInit != nil {
		r.AutoInit = rp.AutoInit
	}
	if rp.IsTemplate != nil {
		r.IsTemplate = rp.IsTemplate
	}
	if rp.TeamID != nil {
		r.TeamID = rp.TeamID
	}
	if rp.GitignoreTemplate != nil {
		r.GitignoreTemplate = rp.GitignoreTemplate
	}
	if rp.LicenseTemplate != nil {
		r.LicenseTemplate = rp.LicenseTemplate
	}
	if rp.AllowSquashMerge != nil {
		r.AllowSquashMerge = rp.AllowSquashMerge
	}
	if rp.AllowMergeCommit != nil {
		r.AllowMergeCommit = rp.AllowMergeCommit
	}
	if rp.AllowRebaseMerge != nil {
		r.AllowRebaseMerge = rp.AllowRebaseMerge
	}
	if rp.DeleteBranchOnMerge != nil {
		r.DeleteBranchOnMerge = rp.DeleteBranchOnMerge
	}
	if rp.HasPages != nil {
		r.HasPages = rp.HasPages
	}
	if rp.HasDownloads != nil {
		r.HasDownloads = rp.HasDownloads
	}
	if rp.DefaultBranch != nil {
		r.DefaultBranch = rp.DefaultBranch
	}
	if rp.Archived != nil {
		r.Archived = rp.Archived
	}
}

// GenerateObservation produces RepositoryObservation object from *github.Repository object.
func GenerateObservation(r github.Repository) v1alpha1.RepositoryObservation {
	o := v1alpha1.RepositoryObservation{
		ID:               r.ID,
		NodeID:           r.NodeID,
		FullName:         r.FullName,
		Name:             r.Name,
		URL:              r.URL,
		ArchiveURL:       r.ArchiveURL,
		AssigneesURL:     r.AssigneesURL,
		BlobsURL:         r.BlobsURL,
		CollaboratorsURL: r.CollaboratorsURL,
		CommentsURL:      r.CommentsURL,
		CommitsURL:       r.CommitsURL,
		CompareURL:       r.CompareURL,
		ContentsURL:      r.ContentsURL,
		ContributorsURL:  r.ContributorsURL,
		DeploymentsURL:   r.DeploymentsURL,
		DownloadsURL:     r.DownloadsURL,
		EventsURL:        r.EventsURL,
		ForksURL:         r.ForksURL,
		GitCommitsURL:    r.GitCommitsURL,
		GitRefsURL:       r.GitRefsURL,
		GitTagsURL:       r.GitTagsURL,
		HooksURL:         r.HooksURL,
		IssueCommentURL:  r.IssueCommentURL,
		IssueEventsURL:   r.IssueEventsURL,
		IssuesURL:        r.IssuesURL,
		KeysURL:          r.KeysURL,
		LabelsURL:        r.LabelsURL,
		LanguagesURL:     r.LanguagesURL,
		MergesURL:        r.MergesURL,
		MilestonesURL:    r.MilestonesURL,
		NotificationsURL: r.NotificationsURL,
		PullsURL:         r.PullsURL,
		ReleasesURL:      r.ReleasesURL,
		StargazersURL:    r.StargazersURL,
		StatusesURL:      r.StatusesURL,
		SubscribersURL:   r.SubscribersURL,
		SubscriptionURL:  r.SubscriptionURL,
		TagsURL:          r.TagsURL,
		TreesURL:         r.TreesURL,
		TeamsURL:         r.TeamsURL,
		HTMLURL:          r.HTMLURL,
		CloneURL:         r.CloneURL,
		GitURL:           r.GitURL,
		MirrorURL:        r.MirrorURL,
		SSHURL:           r.SSHURL,
		SVNURL:           r.SVNURL,
		ForksCount:       r.ForksCount,
		NetworkCount:     r.NetworkCount,
		OpenIssuesCount:  r.OpenIssuesCount,
		StargazersCount:  r.StargazersCount,
		SubscribersCount: r.SubscribersCount,
		WatchersCount:    r.WatchersCount,
		CreatedAt:        ghclient.TimestampConverter(r.CreatedAt),
		PushedAt:         ghclient.TimestampConverter(r.PushedAt),
		UpdatedAt:        ghclient.TimestampConverter(r.UpdatedAt),
		Language:         r.Language,
		Fork:             r.Fork,
		Size:             r.Size,
		Disabled:         r.Disabled,
		Topics:           r.Topics,
		Permissions:      *r.Permissions,
	}
	return o
}

// LateInitialize fills the empty fields of RepositoryParameters if the corresponding
// fields are given in Repository.
func LateInitialize(rp *v1alpha1.RepositoryParameters, r github.Repository) { // nolint:gocyclo
	if rp.Organization == nil && r.Organization.Login != nil {
		rp.Organization = r.Organization.Login
	}
	if rp.Description == nil && r.Description != nil {
		rp.Description = r.Description
	}
	if rp.Homepage == nil && r.Homepage != nil {
		rp.Homepage = r.Homepage
	}
	if rp.Private == nil && r.Private != nil {
		rp.Private = r.Private
	}
	if rp.Visibility == nil && r.Visibility != nil {
		rp.Visibility = r.Visibility
	}
	if rp.HasIssues == nil && r.HasIssues != nil {
		rp.HasIssues = r.HasIssues
	}
	if rp.HasProjects == nil && r.HasProjects != nil {
		rp.HasProjects = r.HasProjects
	}
	if rp.HasWiki == nil && r.HasWiki != nil {
		rp.HasWiki = r.HasWiki
	}
	if rp.IsTemplate == nil && r.IsTemplate != nil {
		rp.IsTemplate = r.IsTemplate
	}
	if rp.TeamID == nil && r.TeamID != nil {
		rp.TeamID = r.TeamID
	}
	if rp.AutoInit == nil && r.AutoInit != nil {
		rp.AutoInit = r.AutoInit
	}
	if rp.GitignoreTemplate == nil && r.GitignoreTemplate != nil {
		rp.GitignoreTemplate = r.GitignoreTemplate
	}
	if rp.LicenseTemplate == nil && r.LicenseTemplate != nil {
		rp.LicenseTemplate = r.LicenseTemplate
	}
	if rp.AllowSquashMerge == nil && r.AllowSquashMerge != nil {
		rp.AllowSquashMerge = r.AllowSquashMerge
	}
	if rp.AllowMergeCommit == nil && r.AllowMergeCommit != nil {
		rp.AllowMergeCommit = r.AllowMergeCommit
	}
	if rp.AllowRebaseMerge == nil && r.AllowRebaseMerge != nil {
		rp.AllowRebaseMerge = r.AllowRebaseMerge
	}
	if rp.DeleteBranchOnMerge == nil && r.DeleteBranchOnMerge != nil {
		rp.DeleteBranchOnMerge = r.DeleteBranchOnMerge
	}
	if rp.HasPages == nil && r.HasPages != nil {
		rp.HasPages = r.HasPages
	}
	if rp.HasDownloads == nil && r.HasDownloads != nil {
		rp.HasDownloads = r.HasDownloads
	}
	if rp.DefaultBranch == nil && r.DefaultBranch != nil {
		rp.DefaultBranch = r.DefaultBranch
	}
	if rp.Archived == nil && r.Archived != nil {
		rp.Archived = r.Archived
	}
}

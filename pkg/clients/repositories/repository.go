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
	"context"

	"strings"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/go-github/v33/github"
	"github.com/mitchellh/copystructure"
	"github.com/pkg/errors"
)

const (
	errCheckUpToDate = "unable to determine if external resource is up to date"
	errFullname      = "The templateRef fullname is not valid. It needs to be in the format {owner}/{name}"
)

// Service defines the Repositories operations
type Service interface {
	Create(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error)
	Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
	Edit(ctx context.Context, owner, repo string, repository *github.Repository) (*github.Repository, *github.Response, error)
	Delete(ctx context.Context, owner, repo string) (*github.Response, error)
	CreateFromTemplate(ctx context.Context, templateOwner, templateRepo string, templateRepoReq *github.TemplateRepoRequest) (*github.Repository, *github.Response, error)
}

// NewService creates a new Service based on the *github.Client
// returned by the NewClient SDK method.
func NewService(token string) *Service {
	c := ghclient.NewClient(token)
	r := Service(c.Repositories)
	return &r
}

// IsUpToDate checks whether Repository is configured with given RepositoryParameters.
func IsUpToDate(rp *v1alpha1.RepositoryParameters, observed *github.Repository) (bool, error) {
	generated, err := copystructure.Copy(observed)
	if err != nil {
		return true, errors.Wrap(err, errCheckUpToDate)
	}
	clone, ok := generated.(*github.Repository)
	if !ok {
		return true, errors.New(errCheckUpToDate)
	}

	desired := OverrideParameters(*rp, *clone)

	return cmp.Equal(
		desired,
		*observed,
		cmpopts.IgnoreFields(github.Repository{}, "AutoInit"),
	), nil
}

// OverrideParameters override the parameters in github.Repository
// that are defined in RepositoryParameters
func OverrideParameters(rp v1alpha1.RepositoryParameters, r github.Repository) github.Repository { // nolint:gocyclo
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
	return r
}

// GenerateObservation produces RepositoryObservation object from github.Repository object.
func GenerateObservation(r github.Repository) v1alpha1.RepositoryObservation {
	o := v1alpha1.RepositoryObservation{
		ID:               ghclient.Int64Value(r.ID),
		NodeID:           ghclient.StringValue(r.NodeID),
		FullName:         ghclient.StringValue(r.FullName),
		Name:             ghclient.StringValue(r.Name),
		URL:              ghclient.StringValue(r.URL),
		ArchiveURL:       ghclient.StringValue(r.ArchiveURL),
		AssigneesURL:     ghclient.StringValue(r.AssigneesURL),
		BlobsURL:         ghclient.StringValue(r.BlobsURL),
		CollaboratorsURL: ghclient.StringValue(r.CollaboratorsURL),
		CommentsURL:      ghclient.StringValue(r.CommentsURL),
		CommitsURL:       ghclient.StringValue(r.CommitsURL),
		CompareURL:       ghclient.StringValue(r.CompareURL),
		ContentsURL:      ghclient.StringValue(r.ContentsURL),
		ContributorsURL:  ghclient.StringValue(r.ContributorsURL),
		DeploymentsURL:   ghclient.StringValue(r.DeploymentsURL),
		DownloadsURL:     ghclient.StringValue(r.DownloadsURL),
		EventsURL:        ghclient.StringValue(r.EventsURL),
		ForksURL:         ghclient.StringValue(r.ForksURL),
		GitCommitsURL:    ghclient.StringValue(r.GitCommitsURL),
		GitRefsURL:       ghclient.StringValue(r.GitRefsURL),
		GitTagsURL:       ghclient.StringValue(r.GitTagsURL),
		HooksURL:         ghclient.StringValue(r.HooksURL),
		IssueCommentURL:  ghclient.StringValue(r.IssueCommentURL),
		IssueEventsURL:   ghclient.StringValue(r.IssueEventsURL),
		IssuesURL:        ghclient.StringValue(r.IssuesURL),
		KeysURL:          ghclient.StringValue(r.KeysURL),
		LabelsURL:        ghclient.StringValue(r.LabelsURL),
		LanguagesURL:     ghclient.StringValue(r.LanguagesURL),
		MergesURL:        ghclient.StringValue(r.MergesURL),
		MilestonesURL:    ghclient.StringValue(r.MilestonesURL),
		NotificationsURL: ghclient.StringValue(r.NotificationsURL),
		PullsURL:         ghclient.StringValue(r.PullsURL),
		ReleasesURL:      ghclient.StringValue(r.ReleasesURL),
		StargazersURL:    ghclient.StringValue(r.StargazersURL),
		StatusesURL:      ghclient.StringValue(r.StatusesURL),
		SubscribersURL:   ghclient.StringValue(r.SubscribersURL),
		SubscriptionURL:  ghclient.StringValue(r.SubscriptionURL),
		TagsURL:          ghclient.StringValue(r.TagsURL),
		TreesURL:         ghclient.StringValue(r.TreesURL),
		TeamsURL:         ghclient.StringValue(r.TeamsURL),
		HTMLURL:          ghclient.StringValue(r.HTMLURL),
		CloneURL:         ghclient.StringValue(r.CloneURL),
		GitURL:           ghclient.StringValue(r.GitURL),
		MirrorURL:        ghclient.StringValue(r.MirrorURL),
		SSHURL:           ghclient.StringValue(r.SSHURL),
		SVNURL:           ghclient.StringValue(r.SVNURL),
		ForksCount:       ghclient.IntValue(r.ForksCount),
		NetworkCount:     ghclient.IntValue(r.NetworkCount),
		OpenIssuesCount:  ghclient.IntValue(r.OpenIssuesCount),
		StargazersCount:  ghclient.IntValue(r.StargazersCount),
		SubscribersCount: ghclient.IntValue(r.SubscribersCount),
		WatchersCount:    ghclient.IntValue(r.WatchersCount),
		CreatedAt:        ghclient.ConvertTimestamp(r.CreatedAt),
		PushedAt:         ghclient.ConvertTimestamp(r.PushedAt),
		UpdatedAt:        ghclient.ConvertTimestamp(r.UpdatedAt),
		Language:         ghclient.StringValue(r.Language),
		Fork:             ghclient.BoolValue(r.Fork),
		Size:             ghclient.IntValue(r.Size),
		Disabled:         ghclient.BoolValue(r.Disabled),
		Topics:           r.Topics,
	}
	if r.Permissions != nil {
		o.Permissions = map[string]bool{}
		for k, v := range *r.Permissions {
			o.Permissions[k] = v
		}
	}

	return o
}

// LateInitialize fills the empty fields of RepositoryParameters if the corresponding
// fields are given in Repository.
func LateInitialize(rp *v1alpha1.RepositoryParameters, r *github.Repository, c xpv1.Condition) { // nolint:gocyclo
	if rp.Organization == nil && ghclient.StringValue(r.Owner.Type) == "Organization" {
		if r.Organization.Login != nil {
			rp.Organization = r.Organization.Login
		}
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
	if rp.Archived == nil && r.Archived != nil {
		rp.Archived = r.Archived
	}
	if r.TemplateRepository != nil {
		rp.Template = &xpv1.Reference{
			Name: *r.TemplateRepository.FullName,
		}
	}

	// This condition below is necessary because the GitHub API is not strongly
	// consistent. In the first moments of creating a repository based on a
	// template, the API always returns that the default branch is "main", then
	// it is changed to the same default branch as the template. This can cause
	// inconsistency in the desired state of the Repository (importing the wrong
	// value into the default branch).
	if c.Reason == xpv1.ReasonCreating && r.TemplateRepository != nil {
		if rp.DefaultBranch == nil {
			rp.DefaultBranch = r.TemplateRepository.DefaultBranch

			// We change the r.DefaultBranch to have consistency when checking
			// if the repository is up to date.
			r.DefaultBranch = r.TemplateRepository.DefaultBranch
		}
	} else {
		if rp.DefaultBranch == nil && r.DefaultBranch != nil {
			rp.DefaultBranch = r.DefaultBranch
		}
	}
}

// SplitFullName splits the repository fullname into map[string]string
// with the keys being "owner" and "name".
func SplitFullName(fullname string) (map[string]string, error) {
	split := strings.Split(fullname, "/")
	if len(split) != 2 {
		return nil, errors.New(errFullname)
	}

	return map[string]string{
		"owner": split[0],
		"name":  split[1],
	}, nil
}

// GenerateTemplateRepoRequest overrides the parameters in github.TemplateRepoRequest
// that are defined in RepositoryParameters.
func GenerateTemplateRepoRequest(rp v1alpha1.RepositoryParameters) github.TemplateRepoRequest {
	r := github.TemplateRepoRequest{}
	if len(rp.Name) != 0 {
		r.Name = ghclient.StringPtr(rp.Name)
	}
	if len(rp.Owner) != 0 {
		r.Owner = ghclient.StringPtr(rp.Owner)
	}
	if rp.Description != nil {
		r.Description = rp.Description
	}
	if rp.Private != nil {
		r.Private = rp.Private
	}
	return r
}

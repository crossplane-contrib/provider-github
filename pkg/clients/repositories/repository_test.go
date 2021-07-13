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
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v33/github"
	"github.com/pkg/errors"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
)

var (
	fakeID       = int64(1)
	fakeNodeID   = "fAKe"
	fakeFullName = "owner/sample"
	fakeType     = "User"

	name               = "sample"
	fakeOwner          = "owner"
	description        = "sample description"
	fakeHasIssues      = false
	fakePrivate        = false
	fakeHasProjects    = false
	fakeHasWiki        = true
	fakeIsTemplate     = true
	fakeAutoInit       = true
	fakeHasPages       = false
	fakeHasDownloads   = true
	fakeDefaultBranch  = "sample"
	fakeTemplateBranch = "template"
	fakeArchived       = false
	fakeFalse          = false
)

func params() *v1alpha1.RepositoryParameters {
	return &v1alpha1.RepositoryParameters{
		Name:          name,
		Owner:         fakeOwner,
		Description:   &description,
		HasIssues:     &fakeHasIssues,
		Private:       &fakePrivate,
		HasProjects:   &fakeHasProjects,
		HasWiki:       &fakeHasWiki,
		IsTemplate:    &fakeIsTemplate,
		AutoInit:      &fakeAutoInit,
		HasPages:      &fakeHasPages,
		HasDownloads:  &fakeHasDownloads,
		DefaultBranch: &fakeDefaultBranch,
		Archived:      &fakeArchived,
	}
}

func unsyncedRepository() *github.Repository {
	return &github.Repository{
		ID:       &fakeID,
		NodeID:   &fakeNodeID,
		FullName: &fakeFullName,
		Owner: &github.User{
			Type: &fakeType,
		},
		Name:          &name,
		Description:   &description,
		HasIssues:     &fakeHasIssues,
		Private:       &fakePrivate,
		HasProjects:   &fakeHasProjects,
		HasWiki:       &fakeFalse,
		IsTemplate:    &fakeFalse,
		AutoInit:      &fakeAutoInit,
		HasPages:      &fakeHasPages,
		HasDownloads:  &fakeHasDownloads,
		DefaultBranch: &fakeDefaultBranch,
		Archived:      &fakeArchived,
	}
}

func syncedRepository() *github.Repository {
	return &github.Repository{
		ID:       &fakeID,
		NodeID:   &fakeNodeID,
		FullName: &fakeFullName,
		Owner: &github.User{
			Type: &fakeType,
		},
		Name:          &params().Name,
		Description:   params().Description,
		HasIssues:     params().HasIssues,
		Private:       params().Private,
		HasProjects:   params().HasProjects,
		HasWiki:       params().HasWiki,
		IsTemplate:    params().IsTemplate,
		AutoInit:      params().AutoInit,
		HasPages:      params().HasPages,
		HasDownloads:  params().HasDownloads,
		DefaultBranch: params().DefaultBranch,
		Archived:      params().Archived,
	}
}

func TestOverrideParameters(t *testing.T) {
	type args struct {
		repo *github.Repository
		rp   v1alpha1.RepositoryParameters
	}
	cases := map[string]struct {
		args
		out *github.Repository
	}{
		"Must create a *github.Repository from RepositoryParameters": {
			args: args{
				rp:   *params(),
				repo: unsyncedRepository(),
			},
			out: syncedRepository(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := OverrideParameters(tc.args.rp, *tc.args.repo)
			if diff := cmp.Diff(*tc.out, got); diff != "" {
				t.Errorf("OverrideParameters(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestLateInitialize(t *testing.T) {
	type args struct {
		repo   github.Repository
		rp     *v1alpha1.RepositoryParameters
		reason xpv1.Condition
	}
	cases := map[string]struct {
		args
		out *v1alpha1.RepositoryParameters
	}{
		"Must use template fields in initialization": {
			args: args{
				repo: github.Repository{
					Name: &params().Name,
					Owner: &github.User{
						Type: &fakeType,
					},
					DefaultBranch: &fakeDefaultBranch,
					TemplateRepository: &github.Repository{
						DefaultBranch: &fakeTemplateBranch,
						FullName:      &fakeFullName,
					},
				},
				rp: &v1alpha1.RepositoryParameters{
					Name:  params().Name,
					Owner: params().Owner,
					Template: &xpv1.Reference{
						Name: fakeFullName,
					},
				},
				reason: xpv1.Condition{
					Reason: xpv1.ReasonCreating,
				},
			},
			out: &v1alpha1.RepositoryParameters{
				Name:  params().Name,
				Owner: params().Owner,
				Template: &xpv1.Reference{
					Name: fakeFullName,
				},
				DefaultBranch: &fakeTemplateBranch,
			},
		},
		"Must initialize template spec field in initialization": {
			args: args{
				repo: github.Repository{
					Name: &params().Name,
					Owner: &github.User{
						Type: &fakeType,
					},
					DefaultBranch: &fakeDefaultBranch,
					TemplateRepository: &github.Repository{
						DefaultBranch: &fakeTemplateBranch,
						FullName:      &fakeFullName,
					},
				},
				rp: &v1alpha1.RepositoryParameters{
					Name:  params().Name,
					Owner: params().Owner,
				},
				reason: xpv1.Condition{},
			},
			out: &v1alpha1.RepositoryParameters{
				Name:  params().Name,
				Owner: params().Owner,
				Template: &xpv1.Reference{
					Name: fakeFullName,
				},
				DefaultBranch: &fakeDefaultBranch,
			},
		},
		"Must initialize empty RepositoryParameters fields if they are given in github.Repository": {
			args: args{
				repo: *syncedRepository(),
				rp: &v1alpha1.RepositoryParameters{
					Name:  params().Name,
					Owner: params().Owner,
				},
				reason: xpv1.Condition{},
			},
			out: params(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			LateInitialize(tc.args.rp, &tc.args.repo, tc.reason)
			if diff := cmp.Diff(tc.args.rp, tc.out); diff != "" {
				t.Errorf("LateInitialize(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestIsUpToDate(t *testing.T) {
	type args struct {
		repo *github.Repository
		rp   *v1alpha1.RepositoryParameters
	}
	cases := map[string]struct {
		args
		out bool
		err error
	}{
		"NotUpToDate": {
			args: args{
				repo: unsyncedRepository(),
				rp:   params(),
			},
			out: false,
		},
		"UpToDate": {
			args: args{
				repo: syncedRepository(),
				rp:   params(),
			},
			out: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, _ := IsUpToDate(tc.args.rp, tc.args.repo)
			if diff := cmp.Diff(tc.out, got); diff != "" {
				t.Errorf("IsUpToDate(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateObservation(t *testing.T) {
	type args struct {
		repo github.Repository
	}
	cases := map[string]struct {
		args
		out v1alpha1.RepositoryObservation
	}{
		"Must generate an RepositoryObservation based on the given github.Repository": {
			args: args{
				repo: *syncedRepository(),
			},
			out: v1alpha1.RepositoryObservation{
				ID:       *syncedRepository().ID,
				NodeID:   *syncedRepository().NodeID,
				FullName: *syncedRepository().FullName,
				Name:     *syncedRepository().Name,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := GenerateObservation(tc.args.repo)
			if diff := cmp.Diff(tc.out, got); diff != "" {
				t.Errorf("GenerateObservation(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestSplitFullName(t *testing.T) {
	type args struct {
		fullname string
	}
	type want struct {
		repo map[string]string
		err  error
	}
	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"SplitSuccessfull": {
			reason: "Must split the fullname successfully",
			args: args{
				fullname: "crossplane/provider-github",
			},
			want: want{
				repo: map[string]string{
					"owner": "crossplane",
					"name":  "provider-github",
				},
				err: nil,
			},
		},
		"SplitError": {
			reason: "Must return error if fullname is not valid",
			args: args{
				fullname: "crossplane",
			},
			want: want{
				repo: nil,
				err:  errors.New(errFullname),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := SplitFullName(tc.args.fullname)
			if diff := cmp.Diff(tc.want.repo, got); diff != "" {
				t.Errorf("SplitFullName(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("SplitFullName(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestGenerateTemplateRepoRequest(t *testing.T) {
	type args struct {
		rp v1alpha1.RepositoryParameters
	}
	type want struct {
		repo github.TemplateRepoRequest
	}
	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"GenerateSuccessfull": {
			reason: "Must create a github.TemplateRepoRequest from RepositoryParameters",
			args: args{
				rp: *params(),
			},
			want: want{
				repo: github.TemplateRepoRequest{
					Name:        &name,
					Owner:       &fakeOwner,
					Description: &description,
					Private:     &fakePrivate,
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := GenerateTemplateRepoRequest(tc.args.rp)
			if diff := cmp.Diff(tc.want.repo, got); diff != "" {
				t.Errorf("GenerateTemplateRepoRequest(...): -want, +got:\n%s", diff)
			}
		})
	}
}

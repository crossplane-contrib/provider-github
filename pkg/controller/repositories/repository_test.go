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
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	"github.com/crossplane-contrib/provider-github/pkg/clients/repositories"
	"github.com/crossplane-contrib/provider-github/pkg/controller/repositories/fake"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-github/v33/github"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	unexpectedObject resource.Managed
	errBoom          = errors.New("boom")
	notFound         = 404
	ok               = 200
	internalError    = 500
	fakeTrue         = true
	fakeFalse        = false
	fakeType         = "User"
	fakeOwner        = "crossplane"
	fakeName         = "sample"
)

type repositoryOption func(*v1alpha1.Repository)

func newRepository(opts ...repositoryOption) *v1alpha1.Repository {
	r := &v1alpha1.Repository{}

	for _, f := range opts {
		f(r)
	}
	return r
}

func withIssues(issues bool) repositoryOption {
	return func(i *v1alpha1.Repository) { i.Spec.ForProvider.HasIssues = &issues }
}

type args struct {
	kube   client.Client
	mg     resource.Managed
	github repositories.Service
}

func TestObserve(t *testing.T) {
	type want struct {
		eo  managed.ExternalObservation
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ResourceIsNotRepository": {
			reason: "Must return an error resource is not Repository",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"CannotGetRepository": {
			reason: "Must return an error if GET repository fails and the error is different than 404",
			args: args{
				mg: newRepository(),
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{
								Response: &http.Response{
									StatusCode: internalError,
								},
							},
							errBoom
					},
				},
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: errors.Wrap(errBoom, errGetRepository),
			},
		},
		"MustNotReturnError404": {
			reason: "Must not return an error if GET repository returns 404 status code",
			args: args{
				mg: newRepository(),
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{
								Response: &http.Response{
									StatusCode: notFound,
								},
							},
							errBoom
					},
				},
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: nil,
			},
		},
		"forProviderUpdateFailed": {
			reason: "Must return an error if forProvider update fails",
			args: args{
				kube: &test.MockClient{
					MockUpdate: func(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
						return errBoom
					},
				},
				mg: newRepository(),
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						userType := "Organization"
						return &github.Repository{
								Name: &repo,
								Organization: &github.Organization{
									Login: &owner,
								},
								Owner: &github.User{
									Type: &userType,
								},
							},
							&github.Response{},
							nil
					},
				},
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: errors.Wrap(errBoom, errKubeUpdateRepository),
			},
		},
		"forProviderUpdateSuccessful": {
			reason: "Must fill the forProvider struct if it is empty",
			args: args{
				kube: &test.MockClient{
					MockUpdate: func(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
						return nil
					},
				},
				mg: newRepository(),
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						userType := fakeType
						return &github.Repository{
								HasIssues: &fakeFalse,
								Owner: &github.User{
									Type: &userType,
								},
							},
							&github.Response{},
							nil
					},
				},
			},
			want: want{
				eo: managed.ExternalObservation{
					ResourceExists:          true,
					ResourceUpToDate:        true,
					ResourceLateInitialized: true,
				},
				err: nil,
			},
		},
		"RepositoryIsNotUpToDate": {
			reason: "Must return ResourceUpToDate as false if Repository is outdated",
			args: args{
				mg: newRepository(
					withIssues(fakeTrue),
				),
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						userType := fakeType
						return &github.Repository{
								HasIssues: &fakeFalse,
								Owner: &github.User{
									Type: &userType,
								},
							},
							&github.Response{},
							nil
					},
				},
			},
			want: want{
				eo: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: false,
				},
				err: nil,
			},
		},
		"RepositoryIsUpToDate": {
			reason: "Must return ResourceUpToDate as false if Repository is outdated",
			args: args{
				mg: newRepository(
					withIssues(fakeFalse),
				),
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						userType := fakeType
						return &github.Repository{
								HasIssues: &fakeFalse,
								Owner: &github.User{
									Type: &userType,
								},
							},
							&github.Response{},
							nil
					},
				},
			},
			want: want{
				eo: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: true,
				},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{
				client: tc.args.kube,
				gh:     tc.args.github,
			}
			got, err := e.Observe(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.eo, got); diff != "" {
				t.Errorf("Observe(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Observe(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type want struct {
		eo  managed.ExternalCreation
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ResourceIsNotRepository": {
			reason: "Must return an error resource is not Repository",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"CreationFailed": {
			reason: "Must return an error if the repository creation fails",
			args: args{
				github: &fake.MockService{
					MockCreate: func(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{},
							errBoom
					},
				},
				mg: newRepository(),
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.Wrap(errBoom, errCreateRepository),
			},
		},
		"Success": {
			reason: "Must not return an error if everything goes well",
			args: args{
				github: &fake.MockService{
					MockCreate: func(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{},
							nil
					},
				},
				mg: newRepository(),
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{
				gh: tc.args.github,
			}
			got, err := e.Create(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.eo, got); diff != "" {
				t.Errorf("Create(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Create(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type want struct {
		eo  managed.ExternalUpdate
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ResourceIsNotRepository": {
			reason: "Must return an error resource is not Repository",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalUpdate{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"CannotGetRepository": {
			reason: "Must return an error if GET repository fails",
			args: args{
				mg: newRepository(),
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{
								Response: &http.Response{
									StatusCode: internalError,
								},
							},
							errBoom
					},
				},
			},
			want: want{
				eo:  managed.ExternalUpdate{},
				err: errors.Wrap(errBoom, errGetRepository),
			},
		},
		"CannotEditRepository": {
			reason: "Must return an error if update repository fails",
			args: args{
				mg: newRepository(),
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						userType := fakeType
						return &github.Repository{
								Owner: &github.User{
									Type: &userType,
								},
							},
							&github.Response{},
							nil
					},
					MockEdit: func(ctx context.Context, owner, repo string, repository *github.Repository) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{
								Response: &http.Response{
									StatusCode: internalError,
								},
							},
							errBoom
					},
				},
			},
			want: want{
				eo:  managed.ExternalUpdate{},
				err: errors.Wrap(errBoom, errUpdateRepository),
			},
		},
		"Success": {
			reason: "Must not return an error if everything goes well",
			args: args{
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						userType := fakeType
						return &github.Repository{
								Owner: &github.User{
									Type: &userType,
								},
							},
							&github.Response{},
							nil
					},
					MockEdit: func(ctx context.Context, owner, repo string, repository *github.Repository) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{},
							nil
					},
				},
				mg: newRepository(),
			},
			want: want{
				eo:  managed.ExternalUpdate{},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{
				gh: tc.args.github,
			}
			got, err := e.Update(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.eo, got); diff != "" {
				t.Errorf("Update(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Update(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type want struct {
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ResourceIsNotRepository": {
			reason: "Must return an error resource is not Repository",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				err: errors.New(errUnexpectedObject),
			},
		},
		"DeleteFailed": {
			reason: "Must return error if DeleteRepository fails",
			args: args{
				mg: newRepository(),
				github: &fake.MockService{
					MockDelete: func(ctx context.Context, owner string, repo string) (*github.Response, error) {
						return &github.Response{},
							errBoom
					},
				},
			},
			want: want{
				err: errors.Wrap(errBoom, errDeleteRepository),
			},
		},
		"Success": {
			reason: "Must not fail if all calls succeed",
			args: args{
				mg: newRepository(),
				github: &fake.MockService{
					MockDelete: func(ctx context.Context, owner string, repo string) (*github.Response, error) {
						return &github.Response{},
							nil
					},
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{
				gh: tc.args.github,
			}
			err := e.Delete(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Delete(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestGetRepository(t *testing.T) {
	type args struct {
		owner      string
		specName   string
		statusName string
		github     repositories.Service
	}
	type want struct {
		repo     *github.Repository
		response *github.Response
		err      error
	}
	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"GetRepositoryFailed": {
			reason: "Must return an error if both get tries fails",
			args: args{
				owner:      "sample",
				specName:   "sample",
				statusName: "sample2",
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{},
							errBoom
					},
				},
			},
			want: want{
				err:      errBoom,
				repo:     nil,
				response: &github.Response{},
			},
		},
		"GetRepositoryWithSpecName": {
			reason: "Must successfully return an repository when using specName in the GET",
			args: args{
				owner:      "sample",
				specName:   "sample",
				statusName: "sample2",
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{},
							nil
					},
				},
			},
			want: want{
				err:      nil,
				repo:     &github.Repository{},
				response: &github.Response{},
			},
		},
		"GetRepositoryWithStatusName": {
			reason: "Must successfully return an repository when using statusName in the GET",
			args: args{
				owner:      "sample",
				specName:   "sample",
				statusName: "sample2",
				github: &fake.MockService{
					MockGet: func(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
						if repo == "sample" {
							return &github.Repository{},
								&github.Response{},
								errBoom
						}
						return &github.Repository{},
							&github.Response{},
							nil
					},
				},
			},
			want: want{
				err:      nil,
				repo:     &github.Repository{},
				response: &github.Response{},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{
				gh: tc.args.github,
			}
			repo, res, err := e.GetRepository(context.Background(), tc.args.owner, tc.args.specName, tc.args.statusName)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("GetRepository(...): -want error, +got error:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.repo, repo); diff != "" {
				t.Errorf("GetRepository(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.response, res); diff != "" {
				t.Errorf("GetRepository(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestCreateRepository(t *testing.T) {
	type args struct {
		rp     v1alpha1.RepositoryParameters
		github repositories.Service
	}
	type want struct {
		err error
	}
	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"CreateRepository": {
			reason: "Must create a normal repository",
			args: args{
				rp: v1alpha1.RepositoryParameters{
					Owner: fakeOwner,
					Name:  fakeName,
				},
				github: &fake.MockService{
					MockCreate: func(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{},
							nil
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		"CreateRepositoryBasedOnTemplate": {
			reason: "Must create a repository based on template",
			args: args{
				rp: v1alpha1.RepositoryParameters{
					Owner: fakeOwner,
					Name:  fakeName,
					Template: &v1.Reference{
						Name: "crossplane/provider-template",
					},
				},
				github: &fake.MockService{
					MockCreateFromTemplate: func(ctx context.Context, templateOwner string, templateRepo string, templateRepoReq *github.TemplateRepoRequest) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{
								Response: &http.Response{
									StatusCode: ok,
								},
							},
							nil
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		"RepositoryTemplateNotFound": {
			reason: "Must fail when creating a repository based on template that doesn't exists",
			args: args{
				rp: v1alpha1.RepositoryParameters{
					Owner: fakeOwner,
					Name:  fakeName,
					Template: &v1.Reference{
						Name: "crossplane/provider-templaet",
					},
				},
				github: &fake.MockService{
					MockCreateFromTemplate: func(ctx context.Context, templateOwner string, templateRepo string, templateRepoReq *github.TemplateRepoRequest) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{
								Response: &http.Response{
									StatusCode: notFound,
								},
							},
							errBoom
					},
				},
			},
			want: want{
				err: errors.Wrap(errBoom, errTemplateNotFound),
			},
		},
		"CreateRepositoryTemplateInternalError": {
			reason: "Must fail when API returns a error when creating a repository based on template",
			args: args{
				rp: v1alpha1.RepositoryParameters{
					Owner: fakeOwner,
					Name:  fakeName,
					Template: &v1.Reference{
						Name: "crossplane/provider-template",
					},
				},
				github: &fake.MockService{
					MockCreateFromTemplate: func(ctx context.Context, templateOwner string, templateRepo string, templateRepoReq *github.TemplateRepoRequest) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{
								Response: &http.Response{
									StatusCode: internalError,
								},
							},
							errBoom
					},
				},
			},
			want: want{
				err: errBoom,
			},
		},
		"CreateRepositoryError": {
			reason: "Must fail when API returns a error when creating a repository",
			args: args{
				rp: v1alpha1.RepositoryParameters{
					Owner: fakeOwner,
					Name:  fakeName,
				},
				github: &fake.MockService{
					MockCreate: func(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{},
							errBoom
					},
				},
			},
			want: want{
				err: errBoom,
			},
		},
		"FailCreateRepositoryWithInvalidTemplateRef": {
			reason: "Must fail templateRef is not valid",
			args: args{
				rp: v1alpha1.RepositoryParameters{
					Owner: fakeOwner,
					Name:  fakeName,
					Template: &v1.Reference{
						Name: "crossplane",
					},
				},
				github: &fake.MockService{
					MockCreateFromTemplate: func(ctx context.Context, templateOwner string, templateRepo string, templateRepoReq *github.TemplateRepoRequest) (*github.Repository, *github.Response, error) {
						return &github.Repository{},
							&github.Response{},
							nil
					},
				},
			},
			want: want{
				err: errors.New("The templateRef fullname is not valid. It needs to be in the format {owner}/{name}"),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{
				gh: tc.args.github,
			}
			err := e.CreateRepository(context.Background(), tc.args.rp)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("CreateRepository(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

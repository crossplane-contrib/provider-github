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

package secrets

import (
	"context"
	"errors"
	"testing"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane-contrib/provider-github/apis/secrets/v1alpha1"
	gc "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/crossplane-contrib/provider-github/pkg/fake"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v33/github"
	perr "github.com/pkg/errors"
)

var (
	errBoom           = errors.New("boom")
	errExtractSecret  = perr.New("cannot extract from secret key when none specified")
	fakeHashCorrect   = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	fakeHashFalse     = "fakeHash"
	fakeUpdateTime    = time.Time{}
	lastUpdateFalse   = "anytime"
	LastUpdateCorrect = fakeUpdateTime.String()
	owner             = "fakeOwner"
	repository        = "fakeRepository"
	value             = v1.SecretKeySelector{
		SecretReference: v1.SecretReference{
			Name:      "fakeNameReference",
			Namespace: "fakeNamespaceReference",
		},
		Key: "fakeKeySelector",
	}
)

func mockSpec(value *v1.SecretKeySelector) *v1alpha1.SecretsParameters {
	spec := v1alpha1.SecretsParameters{
		Owner:      owner,
		Repository: repository,
		Value:      value,
	}

	return &spec
}

func mockStatus(hash string, lastUpdate string) *v1alpha1.SecretsObservation {
	status := v1alpha1.SecretsObservation{
		EncryptValue: &hash,
		LastUpdate:   &lastUpdate,
	}

	return &status
}

func TestCreateOrUpdateSec(t *testing.T) {
	type args struct {
		ctx    context.Context
		p      *v1alpha1.SecretsParameters
		name   string
		client client.Client
		gh     Service
	}

	type want struct {
		hash string
		err  error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"CannotCallSecret": {
			reason: "Error to call encrypted secret",
			args: args{
				ctx:    context.Background(),
				p:      mockSpec(&value),
				name:   "fakeTest",
				client: test.NewMockClient(),
				gh: &fake.MockService{
					MockGetRepoPublicKey: func(ctx context.Context, owner, repo string) (*github.PublicKey, *github.Response, error) {
						return &github.PublicKey{KeyID: gc.StringPtr("172354871263548712365487"), Key: gc.StringPtr("ZjRrM2szeQ==")}, &github.Response{}, errBoom
					},
				},
			},
			want: want{
				hash: "",
				err:  errBoom,
			},
		},
		"CannotCreateOrUpdateSecret": {
			reason: "Error to create repository secret in GitHub",
			args: args{
				ctx:    context.Background(),
				p:      mockSpec(&value),
				name:   "fakeTest",
				client: test.NewMockClient(),
				gh: &fake.MockService{
					MockGetRepoPublicKey: func(ctx context.Context, owner, repo string) (*github.PublicKey, *github.Response, error) {
						return &github.PublicKey{KeyID: gc.StringPtr("172354871263548712365487"), Key: gc.StringPtr("ZjRrM2szeQ==")}, &github.Response{}, nil
					},
					MockCreateOrUpdateRepoSecret: func(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) (*github.Response, error) {
						return &github.Response{}, errBoom
					},
				},
			},
			want: want{
				hash: "",
				err:  errBoom,
			},
		},
		"Success": {
			reason: "If is all correct and return error nil",
			args: args{
				ctx:    context.Background(),
				p:      mockSpec(&value),
				name:   "fakeTest",
				client: test.NewMockClient(),
				gh: &fake.MockService{
					MockGetRepoPublicKey: func(ctx context.Context, owner, repo string) (*github.PublicKey, *github.Response, error) {
						return &github.PublicKey{KeyID: gc.StringPtr("172354871263548712365487"), Key: gc.StringPtr("ZjRrM2szeQ==")}, &github.Response{}, nil
					},
					MockCreateOrUpdateRepoSecret: func(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) (*github.Response, error) {
						return &github.Response{}, nil
					},
				},
			},
			want: want{
				hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				err:  nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := CreateOrUpdateSec(tc.args.ctx, tc.args.p, tc.args.name, tc.args.client, tc.args.gh)
			if diff := cmp.Diff(tc.want.hash, got); diff != "" {
				t.Errorf("CreateOrUpdateSec(...): -want, +got:\n%s", diff)
			}

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("CreateOrUpdateSec(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestIsUpToDate(t *testing.T) {
	type args struct {
		ctx    context.Context
		client client.Client
		p      *v1alpha1.SecretsParameters
		o      *v1alpha1.SecretsObservation
		name   string
		gh     Service
	}

	type want struct {
		ud  bool
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"CannotGetRepoSecret": {
			reason: "Error to get encrypted secret in GitHub",
			args: args{
				ctx:    context.Background(),
				client: test.NewMockClient(),
				p:      mockSpec(&value),
				o:      mockStatus(fakeHashCorrect, LastUpdateCorrect),
				name:   "fakeTest",
				gh: &fake.MockService{
					MockGetRepoSecret: func(ctx context.Context, owner, repo, name string) (*github.Secret, *github.Response, error) {
						return &github.Secret{}, &github.Response{}, errBoom
					},
				},
			},
			want: want{
				ud:  false,
				err: errBoom,
			},
		},
		"CannotExtractSecret": {
			reason: "Error to extract secret in k8s",
			args: args{
				ctx:    context.Background(),
				p:      mockSpec(nil),
				o:      mockStatus(fakeHashCorrect, LastUpdateCorrect),
				name:   "fakeTest",
				client: test.NewMockClient(),
				gh: &fake.MockService{
					MockGetRepoSecret: func(ctx context.Context, owner, repo, name string) (*github.Secret, *github.Response, error) {
						return &github.Secret{Name: "TESTSECRET", CreatedAt: github.Timestamp{Time: fakeUpdateTime}, UpdatedAt: github.Timestamp{Time: fakeUpdateTime}}, &github.Response{}, nil
					},
				},
			},
			want: want{
				ud:  false,
				err: errExtractSecret,
			},
		},
		"HashIsNotUpToDate": {
			reason: "Hash secret is not up to date",
			args: args{
				ctx:    context.Background(),
				p:      mockSpec(&value),
				o:      mockStatus(fakeHashFalse, LastUpdateCorrect),
				name:   "fakeTest",
				client: test.NewMockClient(),
				gh: &fake.MockService{
					MockGetRepoSecret: func(ctx context.Context, owner, repo, name string) (*github.Secret, *github.Response, error) {
						return &github.Secret{Name: "TESTSECRET", CreatedAt: github.Timestamp{Time: fakeUpdateTime}, UpdatedAt: github.Timestamp{Time: fakeUpdateTime}}, &github.Response{}, nil
					},
				},
			},
			want: want{
				ud:  false,
				err: nil,
			},
		},
		"TimeIsNotUpToDate": {
			reason: "Last time updated secret repository is not up to date",
			args: args{
				ctx:    context.Background(),
				p:      mockSpec(&value),
				o:      mockStatus(fakeHashCorrect, lastUpdateFalse),
				name:   "fakeTest",
				client: test.NewMockClient(),
				gh: &fake.MockService{
					MockGetRepoSecret: func(ctx context.Context, owner, repo, name string) (*github.Secret, *github.Response, error) {
						return &github.Secret{Name: "TESTSECRET", CreatedAt: github.Timestamp{Time: fakeUpdateTime}, UpdatedAt: github.Timestamp{Time: fakeUpdateTime}}, &github.Response{}, nil
					},
				},
			},
			want: want{
				ud:  false,
				err: nil,
			},
		},
		"IsUpToDate": {
			reason: "Everything is up to date",
			args: args{
				ctx:    context.Background(),
				p:      mockSpec(&value),
				o:      mockStatus(fakeHashCorrect, LastUpdateCorrect),
				name:   "fakeTest",
				client: test.NewMockClient(),
				gh: &fake.MockService{
					MockGetRepoSecret: func(ctx context.Context, owner, repo, name string) (*github.Secret, *github.Response, error) {
						return &github.Secret{Name: "TESTSECRET", CreatedAt: github.Timestamp{Time: fakeUpdateTime}, UpdatedAt: github.Timestamp{Time: fakeUpdateTime}}, &github.Response{}, nil
					},
				},
			},
			want: want{
				ud:  true,
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := IsUpToDate(tc.args.ctx, tc.args.client, tc.args.p, tc.args.o, tc.args.name, tc.args.gh)
			if diff := cmp.Diff(tc.want.ud, got); diff != "" {
				t.Errorf("IsUpToDate(...): -want, +got:\n%s", diff)
			}

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("IsUpToDate(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

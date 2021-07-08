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

package fake

import (
	"context"

	"github.com/crossplane-contrib/provider-github/pkg/clients/repositories"
	"github.com/google/go-github/v33/github"
)

// This ensures that the mock implements the Service interface
var _ repositories.Service = (*MockService)(nil)

// MockService is a mock implementation of the Service
type MockService struct {
	MockCreate             func(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error)
	MockGet                func(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
	MockEdit               func(ctx context.Context, owner, repo string, repository *github.Repository) (*github.Repository, *github.Response, error)
	MockDelete             func(ctx context.Context, owner, repo string) (*github.Response, error)
	MockCreateFromTemplate func(ctx context.Context, templateOwner, templateRepo string, templateRepoReq *github.TemplateRepoRequest) (*github.Repository, *github.Response, error)
}

// Create is a fake Create SDK method
func (m *MockService) Create(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
	return m.MockCreate(ctx, org, repo)
}

// Get is a fake Get SDK method
func (m *MockService) Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	return m.MockGet(ctx, owner, repo)
}

// Edit is a fake Edit SDK method
func (m *MockService) Edit(ctx context.Context, owner, repo string, repository *github.Repository) (*github.Repository, *github.Response, error) {
	return m.MockEdit(ctx, owner, repo, repository)
}

// Delete is a fake Delete SDK method
func (m *MockService) Delete(ctx context.Context, owner, repo string) (*github.Response, error) {
	return m.MockDelete(ctx, owner, repo)
}

// CreateFromTemplate is a fake CreateFromTemplate SDK method
func (m *MockService) CreateFromTemplate(ctx context.Context, templateOwner, templateRepo string, templateRepoReq *github.TemplateRepoRequest) (*github.Repository, *github.Response, error) {
	return m.MockCreateFromTemplate(ctx, templateOwner, templateRepo, templateRepoReq)
}

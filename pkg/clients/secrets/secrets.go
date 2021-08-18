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
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/crossplane-contrib/provider-github/apis/secrets/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/google/go-github/v33/github"
	"golang.org/x/crypto/nacl/box"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service defines the Secrets operations
type Service interface {
	GetRepoPublicKey(ctx context.Context, owner, repo string) (*github.PublicKey, *github.Response, error)
	GetRepoSecret(ctx context.Context, owner, repo, name string) (*github.Secret, *github.Response, error)
	CreateOrUpdateRepoSecret(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) (*github.Response, error)
	DeleteRepoSecret(ctx context.Context, owner, repo, name string) (*github.Response, error)
}

// NewService creates a new Service based on the *github.Client
// returned by the NewClient SDK method.
func NewService(token string) *Service {
	c := ghclient.NewClient(token)
	r := Service(c.Actions)
	return &r
}

// CreateOrUpdateSec create or update repository secret in GitHub
func CreateOrUpdateSec(ctx context.Context, cr *v1alpha1.SecretsParameters, name string, client client.Client, gh Service) (string, error) {
	encryptedSecret, hash, err := callEncryptSecret(ctx, client, cr, name, gh)
	if err != nil {
		return "", fmt.Errorf("error to call secret: %v", err)
	}

	if _, err := gh.CreateOrUpdateRepoSecret(ctx, cr.Owner, cr.Repository, encryptedSecret); err != nil {
		return "", fmt.Errorf("error to create or update secret: %v", err)
	}

	return hash, nil
}

// IsUpToDate check if encrypted value is up to date
func IsUpToDate(ctx context.Context, client client.Client, p *v1alpha1.SecretsParameters, o *v1alpha1.SecretsObservation, name string, gh Service) (bool, error) {
	sec, _, err := gh.GetRepoSecret(ctx, p.Owner, p.Repository, name)
	if err != nil {
		return false, fmt.Errorf("error to get secret from github: %v", err)
	}

	_, hash, err := callEncryptSecret(ctx, client, p, name, gh)
	if err != nil {
		return false, fmt.Errorf("error to call secret: %v", err)
	}

	if hash != *o.EncryptValue {
		return false, nil
	}

	if sec.UpdatedAt.String() != *o.LastUpdate {
		return false, nil
	}

	return true, nil
}

func callEncryptSecret(ctx context.Context, client client.Client, cr *v1alpha1.SecretsParameters, name string, gh Service) (*github.EncryptedSecret, string, error) {
	publicKey, _, err := gh.GetRepoPublicKey(ctx, cr.Owner, cr.Repository)
	if err != nil {
		return nil, "", err
	}

	ref := xpv1.CommonCredentialSelectors{SecretRef: &cr.Value}
	val, _ := resource.ExtractSecret(ctx, client, ref)
	encryptedSecret, err := encryptSecret(publicKey, name, string(val))
	if err != nil {
		return nil, "", err
	}

	hash := generateHash(string(val))
	return encryptedSecret, hash, nil
}

func generateHash(secretValue string) string {
	h := sha256.Sum256([]byte(secretValue))
	return fmt.Sprintf("%x", h)
}

func encryptSecret(publicKey *github.PublicKey, secretName string, secretValue string) (*github.EncryptedSecret, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, fmt.Errorf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err)
	}

	var boxKey [32]byte
	copy(boxKey[:], decodedPublicKey)
	secretBytes := []byte(secretValue)
	encryptedBytes, err := box.SealAnonymous([]byte{}, secretBytes, &boxKey, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("box.SealAnonymous failed with error %w", err)
	}

	encryptedString := base64.StdEncoding.EncodeToString(encryptedBytes)
	keyID := publicKey.GetKeyID()
	encryptedSecret := &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyID,
		EncryptedValue: encryptedString,
	}
	return encryptedSecret, nil
}

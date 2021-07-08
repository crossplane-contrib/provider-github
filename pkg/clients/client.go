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

package clients

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-github/apis/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetConfig gets the config.
func GetConfig(ctx context.Context, c client.Client, mg resource.Managed) ([]byte, error) {
	pc := &v1beta1.ProviderConfig{}
	if err := c.Get(ctx, types.NamespacedName{Name: mg.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, "cannot get referenced ProviderConfig")
	}

	t := resource.NewProviderConfigUsageTracker(c, &v1beta1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, "cannot track ProviderConfig usage")
	}

	return resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, c, pc.Spec.Credentials.CommonCredentialSelectors)
}

// NewClient creates a new client.
func NewClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

// StringPtr converts the supplied string to a pointer to that string.
func StringPtr(p string) *string { return &p }

// StringValue converts the supplied pointer string to a string.
func StringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// ConvertTimestamp converts *github.Timestamp into *metav1.Time
func ConvertTimestamp(t *github.Timestamp) *metav1.Time {
	if t == nil {
		return nil
	}
	return &metav1.Time{
		Time: t.Time,
	}
}

// Int64Value converts the supplied pointer int64 to a int64.
func Int64Value(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// IntValue converts the supplied pointer int to a int.
func IntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// BoolValue converts the supplied pointer bool to a bool.
func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

/*
Copyright 2020 The Crossplane Authors.

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

package controller

import (
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	"github.com/crossplane-contrib/provider-github/pkg/controller/config"
	"github.com/crossplane-contrib/provider-github/pkg/controller/organizations"
	"github.com/crossplane-contrib/provider-github/pkg/controller/repositories"
	"github.com/crossplane-contrib/provider-github/pkg/controller/secrets"
)

// Setup creates all GitHub controllers with the supplied logger and adds them
// to the supplied manager.
func Setup(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter) error {
	for _, setup := range []func(ctrl.Manager, logging.Logger, workqueue.RateLimiter) error{
		config.Setup,
		organizations.SetupMembership,
		repositories.SetupRepository,
		secrets.SetupSecrets,
	} {
		if err := setup(mgr, l, rl); err != nil {
			return err
		}
	}
	return nil
}

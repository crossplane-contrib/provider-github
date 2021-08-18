package secrets

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-github/apis/secrets/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/crossplane-contrib/provider-github/pkg/clients/secrets"
)

const (
	errUnexpectedObject = "The managed resource is not a Secrets resource"
	errCreateSecrets    = "cannot create Secrets"
	errUpdateSecrets    = "cannot update Secrets"
	errDeleteSecrets    = "cannot delete Secrets"

// 	errKubeUpdateSecrets = "cannot update Secrets custom resource"
// 	errTemplateNotFound  = "the referenced Secrets template was not found"
)

// SetupSecrets adds a controller that reconciles secrets.
func SetupSecrets(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter) error {
	name := managed.ControllerName(v1alpha1.SecretsGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&v1alpha1.Secrets{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.SecretsGroupVersionKind),
			managed.WithExternalConnecter(
				&connector{
					client:      mgr.GetClient(),
					newClientFn: secrets.NewService,
				},
			),
			managed.WithConnectionPublishers(),
			managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
			managed.WithInitializers(managed.NewDefaultProviderConfig(mgr.GetClient())),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

type connector struct {
	client      client.Client
	newClientFn func(string) *secrets.Service
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Secrets)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	cfg, err := ghclient.GetConfig(ctx, c.client, cr)
	if err != nil {
		return nil, err
	}
	return &external{*c.newClientFn(string(cfg)), c.client}, nil
}

type external struct {
	gh     secrets.Service
	client client.Client
}

func (e *external) Observe(ctx context.Context, mgd resource.Managed) (managed.ExternalObservation, error) {
	ud := true
	ex := false
	cr, ok := mgd.(*v1alpha1.Secrets)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}

	cr.Status.SetConditions(xpv1.Available())
	if cr.Status.AtProvider.EncryptValue != nil {
		ex = true
		check, err := secrets.IsUpToDate(ctx, e.client, &cr.Spec.ForProvider, &cr.Status.AtProvider, meta.GetExternalName(cr), e.gh)
		if err != nil {
			return managed.ExternalObservation{}, errors.Wrap(err, "Error to verify if is up to date")
		}

		if !check {
			ud = false
		}
	}

	return managed.ExternalObservation{
		ResourceUpToDate:  ud,
		ResourceExists:    ex,
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (e *external) Create(ctx context.Context, mgd resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mgd.(*v1alpha1.Secrets)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}

	hash, err := secrets.CreateOrUpdateSec(ctx, &cr.Spec.ForProvider, meta.GetExternalName(cr), e.client, e.gh)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateSecrets)
	}

	sec, _, err := e.gh.GetRepoSecret(ctx, cr.Spec.ForProvider.Owner, cr.Spec.ForProvider.Repository, meta.GetExternalName(cr))
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "Error to get secret from GitHub")
	}

	cr.Status.AtProvider.LastUpdate = ghclient.StringPtr(sec.UpdatedAt.String())
	cr.Status.AtProvider.EncryptValue = &hash
	cr.SetConditions(xpv1.Creating())
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mgd resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mgd.(*v1alpha1.Secrets)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}

	hash, err := secrets.CreateOrUpdateSec(ctx, &cr.Spec.ForProvider, meta.GetExternalName(cr), e.client, e.gh)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateSecrets)
	}

	sec, _, err := e.gh.GetRepoSecret(ctx, cr.Spec.ForProvider.Owner, cr.Spec.ForProvider.Repository, meta.GetExternalName(cr))
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "Error to get secret from GitHub")
	}

	cr.Status.AtProvider.LastUpdate = ghclient.StringPtr(sec.UpdatedAt.String())
	cr.Status.AtProvider.EncryptValue = &hash
	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mgd resource.Managed) error {
	cr, ok := mgd.(*v1alpha1.Secrets)
	if !ok {
		return errors.New(errUnexpectedObject)
	}

	_, err := e.gh.DeleteRepoSecret(ctx, cr.Spec.ForProvider.Owner, cr.Spec.ForProvider.Repository, meta.GetExternalName(cr))
	if err != nil {
		return errors.Wrap(err, errDeleteSecrets)
	}

	cr.Status.AtProvider.LastUpdate = nil
	cr.Status.AtProvider.EncryptValue = nil
	cr.SetConditions(xpv1.Deleting())
	return nil
}

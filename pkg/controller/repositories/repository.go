package repositories

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v33/github"
	"github.com/pkg/errors"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/crossplane-contrib/provider-github/pkg/clients/repositories"
)

const (
	errUnexpectedObject     = "The managed resource is not a Repository resource"
	errGetRepository        = "Cannot get GitHub repository"
	errCheckUpToDate        = "unable to determine if external resource is up to date"
	errCreateRepository     = "cannot create Repository"
	errUpdateRepository     = "cannot update Repository"
	errDeleteRepository     = "cannot delete Repository"
	errKubeUpdateRepository = "cannot update Repository custom resource"
)

// SetupRepository adds a controller that reconciles Repositories.
func SetupRepository(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter) error {
	name := managed.ControllerName(v1alpha1.RepositoryGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&v1alpha1.Repository{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.RepositoryGroupVersionKind),
			managed.WithExternalConnecter(&connector{client: mgr.GetClient(), newClientFn: ghclient.NewClient}),
			managed.WithConnectionPublishers(),
			managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
			managed.WithInitializers(managed.NewDefaultProviderConfig(mgr.GetClient())),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

type connector struct {
	client      client.Client
	newClientFn func(string) *github.Client
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Repository)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	cfg, err := ghclient.GetConfig(ctx, c.client, cr)
	if err != nil {
		return nil, err
	}
	return &external{c.newClientFn(string(cfg)), c.client}, nil
}

type external struct {
	gh     *github.Client
	client client.Client
}

func (e *external) Observe(ctx context.Context, mgd resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mgd.(*v1alpha1.Repository)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}

	r, err := e.GetRepository(
		ctx,
		cr.Spec.ForProvider.Owner,
		cr.Spec.ForProvider.Name,
		ghclient.StringValue(cr.Status.AtProvider.Name),
	)
	if err != nil {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	// Import repository if already exists
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	repositories.LateInitialize(&cr.Spec.ForProvider, *r)
	if !cmp.Equal(currentSpec, &cr.Spec.ForProvider) {
		if err := e.client.Update(ctx, cr); err != nil {
			return managed.ExternalObservation{}, errors.Wrap(err, errKubeUpdateRepository)
		}
	}

	cr.Status.SetConditions(xpv1.Available())
	cr.Status.AtProvider = repositories.GenerateObservation(*r)

	upToDate, err := repositories.IsUpToDate(&cr.Spec.ForProvider, r)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errCheckUpToDate)
	}

	return managed.ExternalObservation{
		ResourceUpToDate: upToDate,
		ResourceExists:   true,
	}, nil
}

func (e *external) Create(ctx context.Context, mgd resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mgd.(*v1alpha1.Repository)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}

	repo := &github.Repository{}
	repositories.GenerateRepository(cr.Spec.ForProvider, repo)

	_, _, err := e.gh.Repositories.Create(
		ctx,
		ghclient.StringValue(cr.Spec.ForProvider.Organization),
		repo,
	)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateRepository)
	}

	cr.SetConditions(xpv1.Creating())

	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mgd resource.Managed) (managed.ExternalUpdate, error) { // nolint:gocyclo
	cr, ok := mgd.(*v1alpha1.Repository)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}

	r, err := e.GetRepository(
		ctx,
		cr.Spec.ForProvider.Owner,
		cr.Spec.ForProvider.Name,
		ghclient.StringValue(cr.Status.AtProvider.Name),
	)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errGetRepository)
	}

	repositories.GenerateRepository(cr.Spec.ForProvider, r)

	_, _, err = e.gh.Repositories.Edit(
		ctx,
		cr.Spec.ForProvider.Owner,
		*cr.Status.AtProvider.Name,
		r,
	)
	return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateRepository)
}

func (e *external) Delete(ctx context.Context, mgd resource.Managed) error {
	cr, ok := mgd.(*v1alpha1.Repository)
	if !ok {
		return errors.New(errUnexpectedObject)
	}

	_, err := e.gh.Repositories.Delete(ctx,
		cr.Spec.ForProvider.Owner,
		cr.Spec.ForProvider.Name,
	)
	return errors.Wrap(err, errDeleteRepository)
}

// GetRepository makes API calls to get the Repository.
// If using the Spec name the repository is not found, a second attempt
// is made with the status name. This is useful when updating the Repository name.
func (e *external) GetRepository(ctx context.Context, owner, specName, statusName string) (*github.Repository, error) {
	r, _, err := e.gh.Repositories.Get(ctx, owner, specName)
	if err != nil {
		r, _, err = e.gh.Repositories.Get(ctx, owner, statusName)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	return r, nil
}

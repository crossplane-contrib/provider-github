module github.com/crossplane-contrib/provider-github

go 1.13

require (
	github.com/crossplane/crossplane-runtime v0.16.1
	github.com/crossplane/crossplane-tools v0.0.0-20201201125637-9ddc70edfd0d
	github.com/google/go-github/v33 v33.0.0
	github.com/google/uuid v1.1.4 // indirect
	github.com/pkg/errors v0.9.1
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	sigs.k8s.io/controller-runtime v0.11.0
	sigs.k8s.io/controller-tools v0.8.0
)

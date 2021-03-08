module github.com/crossplane-contrib/provider-github

go 1.13

require (
	github.com/crossplane/crossplane-runtime v0.13.0
	github.com/crossplane/crossplane-tools v0.0.0-20201201125637-9ddc70edfd0d
	github.com/google/go-github/v33 v33.0.0
	github.com/pkg/errors v0.9.1
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200910180754-dd1b699fc489
	golang.org/x/oauth2 v0.0.0-20210112200429-01de73cf58bd
	google.golang.org/grpc/examples v0.0.0-20210304020650-930c79186c99 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
	k8s.io/release v0.7.0 // indirect
	sigs.k8s.io/controller-runtime v0.8.0
	sigs.k8s.io/controller-tools v0.2.4
)

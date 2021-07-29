module github.com/dlorenc/cosigned

go 1.16

require (
	github.com/docker/cli v20.10.0-beta1.0.20201117192004-5cc239616494+incompatible // indirect
	github.com/docker/docker v20.10.0-beta1.0.20201110211921-af34b94a78a1+incompatible // indirect
	github.com/google/go-containerregistry v0.5.1
	github.com/sigstore/cosign v1.0.1-0.20210728181701-5f1f18426dc3
	github.com/sigstore/sigstore v0.0.0-20210722023421-fd3b69438dba
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
	sigs.k8s.io/controller-runtime v0.9.2
)

replace github.com/prometheus/common => github.com/prometheus/common v0.26.0

module github.com/dlorenc/cosigned

go 1.16

require (
	github.com/docker/cli v20.10.0-beta1.0.20201117192004-5cc239616494+incompatible // indirect
	github.com/docker/docker v20.10.0-beta1.0.20201110211921-af34b94a78a1+incompatible // indirect
	github.com/google/go-containerregistry v0.5.1
	github.com/sigstore/cosign v0.6.1-0.20210713005353-82d49dcf3b8b
	github.com/sigstore/rekor v0.2.1-0.20210712122031-1c30d2ff9518 // indirect
	github.com/sigstore/sigstore v0.0.0-20210709190449-2ab5ec881a5f
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.9.2
)

replace github.com/prometheus/common => github.com/prometheus/common v0.26.0

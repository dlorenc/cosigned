package cosigned

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/sigstore/cosign/cmd/cosign/cli"
	"github.com/sigstore/sigstore/pkg/signature"
	corev1 "k8s.io/api/core/v1"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/cosign/pkg/cosign/fulcio"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	log = ctrl.Log.WithName("cosigned")
)

func Signatures(ctx context.Context, img string, key *ecdsa.PublicKey) ([]cosign.SignedPayload, error) {
	ref, err := name.ParseReference(img)
	if err != nil {
		return nil, err
	}

	ecdsaVerifier := &signature.ECDSAVerifier{Key: key, HashAlg: crypto.SHA256}

	return cosign.Verify(ctx, ref, &cosign.CheckOpts{
		Roots:  fulcio.Roots,
		PubKey: ecdsaVerifier,
		Claims: true,
	}, cli.TlogServer())
}

func Config(ctx context.Context, c client.Client) *corev1.ConfigMap {
	obj := &corev1.ConfigMap{}
	if err := c.Get(ctx, client.ObjectKey{
		Namespace: "cosigned-system",
		Name:      "cosigned-config",
	}, obj); err != nil {
		log.Error(err, "getting configmap")
	}
	return obj
}

func Keys(cfg map[string]string) []*ecdsa.PublicKey {
	keys := []*ecdsa.PublicKey{}

	pems := parsePems([]byte(cfg["keys"]))
	for _, p := range pems {
		// TODO check header
		key, err := x509.ParsePKIXPublicKey(p.Bytes)
		if err != nil {
			log.Error(err, "parsing key", "key", p)
		}
		keys = append(keys, key.(*ecdsa.PublicKey))
	}
	return keys
}

func parsePems(b []byte) []*pem.Block {
	p, rest := pem.Decode(b)
	if p == nil {
		return nil
	}
	pems := []*pem.Block{p}

	if rest != nil {
		return append(pems, parsePems(rest)...)
	}
	return pems
}

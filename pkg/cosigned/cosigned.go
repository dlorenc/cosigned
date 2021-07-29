package cosigned

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/cmd/cosign/cli/fulcio"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/sigstore/pkg/signature"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	log = ctrl.Log.WithName("cosigned")
)

func Signatures(ctx context.Context, img string, key *ecdsa.PublicKey) ([]cosign.SignedPayload, error) {
	ref, err := name.ParseReference(img)
	if err != nil {
		return nil, err
	}

	ecdsaVerifier, err := signature.LoadECDSAVerifier(key, crypto.SHA256)
	if err != nil {
		return nil, err
	}

	return cosign.Verify(ctx, ref, &cosign.CheckOpts{
		RootCerts:     fulcio.Roots,
		SigVerifier:   ecdsaVerifier,
		ClaimVerifier: cosign.SimpleClaimVerifier,
	})
}

func Keys(cfg map[string][]byte) []*ecdsa.PublicKey {
	keys := []*ecdsa.PublicKey{}

	pems := parsePems(cfg["cosign.pub"])
	for _, p := range pems {
		// TODO check header
		key, err := x509.ParsePKIXPublicKey(p.Bytes)
		if err != nil {
			log.Error(err, "parsing key", "cosign.pub", p)
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

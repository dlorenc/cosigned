package cosigned

import (
	"context"
	"crypto/ecdsa"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/cosign/pkg/cosign/fulcio"
)

func Signatures(ctx context.Context, img string, key *ecdsa.PublicKey) ([]cosign.SignedPayload, error) {
	ref, err := name.ParseReference(img)
	if err != nil {
		return nil, err
	}
	return cosign.Verify(ctx, ref, cosign.CheckOpts{
		Roots:  fulcio.Roots,
		PubKey: key,
		Claims: true,
	})
}

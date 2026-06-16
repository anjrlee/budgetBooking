// Package googleauth verifies Google Identity Services ID tokens without
// pulling in the full google.golang.org/api dependency tree.
package googleauth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const certsURL = "https://www.googleapis.com/oauth2/v3/certs"

var validIssuers = map[string]bool{
	"accounts.google.com":         true,
	"https://accounts.google.com": true,
}

// Claims holds the user identity fields extracted from a verified Google ID token.
type Claims struct {
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type googleClaims struct {
	jwt.RegisteredClaims
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

var (
	keysMu     sync.Mutex
	keysCache  map[string]*rsa.PublicKey
	keysExpiry time.Time
)

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwkSet struct {
	Keys []jwk `json:"keys"`
}

// Verify validates a Google ID token's signature, issuer, audience, and
// expiry, returning the embedded user identity claims.
func Verify(ctx context.Context, idToken, clientID string) (*Claims, error) {
	claims := &googleClaims{}
	token, err := jwt.ParseWithClaims(idToken, claims, func(t *jwt.Token) (interface{}, error) {
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("id token missing kid header")
		}
		return googleKey(ctx, kid)
	}, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil {
		return nil, fmt.Errorf("invalid google id token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid google id token")
	}

	if !validIssuers[claims.Issuer] {
		return nil, fmt.Errorf("unexpected issuer %q", claims.Issuer)
	}
	if !slices.Contains(claims.Audience, clientID) {
		return nil, errors.New("id token audience does not match this application")
	}
	if claims.Subject == "" {
		return nil, errors.New("id token missing subject")
	}

	return &Claims{
		Subject:       claims.Subject,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		Name:          claims.Name,
		Picture:       claims.Picture,
	}, nil
}

func googleKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	keysMu.Lock()
	defer keysMu.Unlock()

	if time.Now().After(keysExpiry) || keysCache[kid] == nil {
		if err := refreshKeysLocked(ctx); err != nil {
			return nil, err
		}
	}

	key, ok := keysCache[kid]
	if !ok {
		return nil, fmt.Errorf("unknown google signing key %q", kid)
	}
	return key, nil
}

func refreshKeysLocked(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, certsURL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetching google certs: unexpected status %d", resp.StatusCode)
	}

	var set jwkSet
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return err
	}

	cache := make(map[string]*rsa.PublicKey, len(set.Keys))
	for _, k := range set.Keys {
		if k.Kty != "RSA" {
			continue
		}
		nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
		if err != nil {
			continue
		}
		cache[k.Kid] = &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: int(new(big.Int).SetBytes(eBytes).Int64()),
		}
	}

	keysCache = cache
	keysExpiry = time.Now().Add(time.Hour)
	return nil
}

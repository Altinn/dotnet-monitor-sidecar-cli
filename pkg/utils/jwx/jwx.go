package jwx

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// CreateJWTKey creates a JWT key and returns the token, subject and the public key
func CreateJWTKey() (token, subject, key string, err error) {
	uuID := uuid.New()
	subject = uuID.String()
	j := jwt.New()
	j.Set("aud", "https://github.com/dotnet/dotnet-monitor")
	j.Set("iss", "https://github.com/dotnet/dotnet-monitor/generatekey+MonitorApiKey")
	j.Set("sub", subject)

	raw, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return
	}
	jwkKey, err := jwk.New(raw)
	if err != nil {
		return
	}
	signed, err := jwt.Sign(j, jwa.ES384, jwkKey)
	token = string(signed)
	pubKey, err := jwkKey.PublicKey()
	if err != nil {
		return
	}
	b, err := json.Marshal(pubKey)
	if err != nil {
		return
	}
	key = base64.URLEncoding.EncodeToString(b)
	return
}

package jwx

import (
	"encoding/base64"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"testing"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

func TestCreateJWTKey(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Verify token has expected subject",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, gotSubject, key, err := CreateJWTKey()
			if err != nil {
				t.Errorf("CreateJWTKey() error = %v", err)
				return
			}
			b, err := base64.URLEncoding.DecodeString(key)
			if err != nil {
				t.Errorf("Failed to base64 decode string %v", err)
			}
			parsedKey, err := jwk.ParseKey(b)
			if err != nil {
				t.Errorf("Failed to parse key error: %v", err)
				return
			}
			token, err := jwt.Parse([]byte(gotToken), jwt.WithKey(jwa.ES384, parsedKey))
			if err != nil {
				t.Errorf("Failed to parse returned token error = %v", err)
				return
			}
			if err != nil {
				t.Errorf("Failed to parse returned key error = %v", err)
				return
			}
			if err != nil {
				t.Errorf("Failed to validate token error = %v", err)
				return
			}
			err = jwt.Validate(token, jwt.WithSubject(gotSubject))
			if err != nil {
				t.Errorf("Failed to validate token error = %v", err)
				return
			}
		})
	}
}

package jwx

import (
	"testing"

	"github.com/lestrrat-go/jwx/jwt"
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
			gotToken, gotSubject, _, err := CreateJWTKey()
			if err != nil {
				t.Errorf("CreateJWTKey() error = %v", err)
				return
			}
			token, err := jwt.Parse([]byte(gotToken))
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

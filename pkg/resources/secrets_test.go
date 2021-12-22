package resources

import (
	"fmt"
	"regexp"
	"testing"

	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

func TestGenerateSecret(t *testing.T) {
	type args struct {
		namespace string
		subject   string
		key       string
		owner     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Generate secret with values",
			args: args{
				subject:   utilrand.String(5),
				key:       utilrand.String(5),
				owner:     utilrand.String(5),
				namespace: utilrand.String(5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GenerateSecret(tt.args.namespace, tt.args.subject, tt.args.key, tt.args.owner)
			nameRegxp := regexp.MustCompile(fmt.Sprintf("%s[a-zA-Z0-9_.-]{5}", SecretBaseName))
			if !nameRegxp.MatchString(actual.Name) {
				t.Errorf("Secret name %s does not match expected pattern %s", actual.Name, nameRegxp)
			}
			if string(actual.Data[SubjectKey]) != tt.args.subject {
				t.Errorf("Secret key %s does not match expected key %s", string(actual.Data[SubjectKey]), tt.args.key)
			}
			if string(actual.Data[publicKeyKey]) != tt.args.key {
				t.Errorf("Secret publickey %s does not match expected key %s", string(actual.Data[publicKeyKey]), tt.args.key)
			}
			if actual.Labels[SecretLabel] != tt.args.owner {
				t.Errorf("Secret owner %s does not match expected owner %s", actual.Labels[SecretLabel], tt.args.owner)
			}
		})
	}
}

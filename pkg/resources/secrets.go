package resources

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

const (
	// SecretBaseName is the base name of the secret generated by the cli
	SecretBaseName = "dd-monitor-apikey-"
	// SecretLabel is the label used to identify the owner of the secret
	SecretLabel = "dev.local/dd-secret"
	// SubjectKey is the key for the subject in the secret
	SubjectKey = "Authentication__MonitorApiKey__Subject"
	// publicKeyKey is the key for the public key in the secret
	publicKeyKey = "Authentication__MonitorApiKey__PublicKey"
)

// GenerateSecret generates a kubernetes secret with a JWK public-key and subject
func GenerateSecret(namespace, subject, key, owner string) corev1.Secret {
	secretName := fmt.Sprintf("%s%s", SecretBaseName, utilrand.String(5))
	s := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels: map[string]string{
				SecretLabel: owner,
			},
		},
		Type: corev1.SecretTypeOpaque,
	}
	s.Data = map[string][]byte{
		SubjectKey:   []byte(subject),
		publicKeyKey: []byte(key),
	}
	return s
}

package kubernetes

import (
	"context"
	"fmt"

	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/errors"
	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/resources"
	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/utils/jwx"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateJWKSecret creates a secret with a JWK public-key and subject
func (h *Helper) CreateJWKSecret(ctx context.Context, namespace, owner string) (name string, token string, err error) {
	s, err := h.FetchJWKSecret(ctx, namespace, owner)
	if err != nil && errors.IsNotFound(err) {
		token, subject, key, err := jwx.CreateJWTKey()
		s = resources.GenerateSecret(namespace, subject, key, owner)
		_, err = h.Client.CoreV1().Secrets(namespace).Create(ctx, &s, metav1.CreateOptions{})
		return s.Name, token, err
	}
	return s.Name, "Existing secret, token not available", err
}

// RemoveJWKSecret removes secret
func (h *Helper) RemoveJWKSecret(ctx context.Context, namespace, secretname string) error {
	return h.Client.CoreV1().Secrets(namespace).Delete(ctx, secretname, metav1.DeleteOptions{})
}

// FetchJWKSecret fetches secret from kubernetes based on owner
func (h *Helper) FetchJWKSecret(ctx context.Context, namespace, owner string) (corev1.Secret, error) {
	sl, err := h.Client.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", resources.SecretLabel, owner)})
	if err != nil {
		return corev1.Secret{}, err
	}
	if len(sl.Items) == 0 {
		return corev1.Secret{}, fmt.Errorf("resource not found")
	}
	if len(sl.Items) > 1 {
		return corev1.Secret{}, fmt.Errorf("multiple resources found")
	}
	return sl.Items[0], nil
}

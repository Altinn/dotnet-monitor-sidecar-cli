package kubernetes

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/resources"
	"github.com/lestrrat-go/jwx/v2/jwt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	testclient "k8s.io/client-go/kubernetes/fake"
	clienttesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
)

func TestHelper_FetchJWKSecret(t *testing.T) {
	type args struct {
		ctx       context.Context
		namespace string
		owner     string
	}
	tests := []struct {
		name           string
		args           args
		existingSecret []runtime.Object
		want           corev1.Secret
		wantErr        bool
	}{
		{
			name: "Returns secret",
			args: args{
				ctx:       context.Background(),
				namespace: "test",
				owner:     "test",
			},
			existingSecret: []runtime.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dd-monitor-apikey-test",
						Namespace: "test",
						Labels: map[string]string{
							"dev.local/dd-secret": "test",
						},
					},
					Type: corev1.SecretTypeOpaque,
					Data: map[string][]byte{
						"Authentication__MonitorApiKey__Subject":   []byte("test"),
						"Authentication__MonitorApiKey__PublicKey": []byte("test"),
					},
				},
			},
			want: corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dd-monitor-apikey-test",
					Namespace: "test",
					Labels: map[string]string{
						"dev.local/dd-secret": "test",
					},
				},
				Type: corev1.SecretTypeOpaque,
				Data: map[string][]byte{
					"Authentication__MonitorApiKey__Subject":   []byte("test"),
					"Authentication__MonitorApiKey__PublicKey": []byte("test"),
				},
			},
			wantErr: false,
		},
		{
			name: "No secret with label matches",
			existingSecret: []runtime.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dd-monitor-apikey-test",
						Namespace: "test",
						Labels: map[string]string{
							"dev.local/dd-secret": "test",
						},
					},
					Type: corev1.SecretTypeOpaque,
					Data: map[string][]byte{
						"Authentication__MonitorApiKey__Subject":   []byte("test"),
						"Authentication__MonitorApiKey__PublicKey": []byte("test"),
					},
				},
			},
			args: args{
				ctx:       context.Background(),
				namespace: "test",
				owner:     "another",
			},
			want:    corev1.Secret{},
			wantErr: true,
		},
		{
			name:           "No secret defined",
			existingSecret: []runtime.Object{},
			args: args{
				ctx:       context.Background(),
				namespace: "test",
				owner:     "another",
			},
			want:    corev1.Secret{},
			wantErr: true,
		},
		{
			name: "Multiple secret with label matches",
			existingSecret: []runtime.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dd-monitor-apikey-test",
						Namespace: "test",
						Labels: map[string]string{
							"dev.local/dd-secret": "test",
						},
					},
					Type: corev1.SecretTypeOpaque,
					Data: map[string][]byte{
						"Authentication__MonitorApiKey__Subject":   []byte("test"),
						"Authentication__MonitorApiKey__PublicKey": []byte("test"),
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dd-monitor-apikey-test2",
						Namespace: "test",
						Labels: map[string]string{
							"dev.local/dd-secret": "test",
						},
					},
					Type: corev1.SecretTypeOpaque,
					Data: map[string][]byte{
						"Authentication__MonitorApiKey__Subject":   []byte("test2"),
						"Authentication__MonitorApiKey__PublicKey": []byte("test2"),
					},
				},
			},
			args: args{
				ctx:       context.Background(),
				namespace: "test",
				owner:     "test",
			},
			want:    corev1.Secret{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *testclient.Clientset
			c = testclient.NewSimpleClientset(tt.existingSecret...)
			h := &Helper{
				Client: c,
			}
			got, err := h.FetchJWKSecret(tt.args.ctx, tt.args.namespace, tt.args.owner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.FetchJWKSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Helper.FetchJWKSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelper_RemoveJWKSecret(t *testing.T) {
	type args struct {
		ctx        context.Context
		namespace  string
		secretname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx:        context.Background(),
				namespace:  "test",
				secretname: "dd-monitor-apikey-test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			secrets := make(chan *corev1.Secret, 1)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := createFakeClientWithSecretWatcher(
				ctx,
				&cache.ResourceEventHandlerFuncs{
					DeleteFunc: func(obj interface{}) {
						secrets <- obj.(*corev1.Secret)
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dd-monitor-apikey-test",
						Namespace: "test",
					},
					Type: corev1.SecretTypeOpaque,
					Data: map[string][]byte{
						"Authentication__MonitorApiKey__Subject": []byte("test"),
					},
				},
			)
			h := &Helper{
				Client: c,
			}
			if err := h.RemoveJWKSecret(tt.args.ctx, tt.args.namespace, tt.args.secretname); (err != nil) != tt.wantErr {
				t.Errorf("Helper.RemoveJWKSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
			select {
			case sec := <-secrets:
				t.Logf("Deleted secret: %s", sec.Name)
			case <-time.After(wait.ForeverTestTimeout):
				t.Errorf("Timed out waiting for secret to be deleted")
			}
		})
	}
}

func TestHelper_CreateJWKSecret(t *testing.T) {
	type args struct {
		ctx       context.Context
		namespace string
		owner     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				namespace: "test",
				owner:     "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			secrets := make(chan *corev1.Secret, 1)
			c := createFakeClientWithSecretWatcher(ctx, &cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					secrets <- obj.(*corev1.Secret)
				},
			})
			h := &Helper{
				Client: c,
			}
			secretName, token, err := h.CreateJWKSecret(tt.args.ctx, tt.args.namespace, tt.args.owner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.CreateJWKSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			secretRegexp := regexp.MustCompile(fmt.Sprintf("%s[a-zA-Z0-9_.-]{5}", resources.SecretBaseName))
			if !secretRegexp.MatchString(secretName) {
				t.Errorf("Secretname %s did not match regex %s", secretName, secretRegexp)
			}
			select {
			case sec := <-secrets:
				t.Logf("Created secret: %s", sec.Name)
				parsedToken, err := jwt.Parse([]byte(token), jwt.WithVerify(false))
				if err != nil {
					t.Errorf("Failed to parse token: %v", err)
				}
				err = jwt.Validate(parsedToken, jwt.WithSubject(string(sec.Data[resources.SubjectKey])))
				if err != nil {
					t.Errorf("Failed to validate token with subject %s error = %v", sec.Data[resources.SubjectKey], err)
					return
				}
			case <-time.After(wait.ForeverTestTimeout):
				t.Errorf("Timed out waiting for secret to be created")
			}
		})
	}
}

func createFakeClientWithSecretWatcher(ctx context.Context, handlers *cache.ResourceEventHandlerFuncs, objs ...runtime.Object) *testclient.Clientset {
	watcherStarted := make(chan struct{})
	c := testclient.NewSimpleClientset(objs...)
	c.PrependWatchReactor("*", func(action clienttesting.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := c.Tracker().Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		close(watcherStarted)
		return true, watch, nil
	})

	i := informers.NewSharedInformerFactory(c, 0)
	secInformers := i.Core().V1().Secrets().Informer()
	secInformers.AddEventHandler(handlers)
	i.Start(ctx.Done())
	<-watcherStarted
	return c
}

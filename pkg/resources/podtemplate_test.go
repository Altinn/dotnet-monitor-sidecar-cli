package resources

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/sergi/go-diff/diffmatchpatch"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
func TestAddDebugContainerPodTemplate(t *testing.T) {
	type args struct {
		namespace        string
		containerToDebug string
		debugimage       string
		secretname       string
	}
	tests := []struct {
		name       string
		args       args
		inputfile  string
		goldenfile string
		wantErr    bool
	}{
		{
			name: "Add debug container to pod template",
			args: args{
				namespace:        "test",
				containerToDebug: "",
				debugimage:       "test:latest",
				secretname:       "secret",
			},
			inputfile:  "testdata/add-pod-template/podtemplate_test_one_container.yaml",
			goldenfile: "testdata/add-pod-template/podtemplate_test_one_container.golden",
			wantErr:    false,
		},
		{
			name: "Add debug container to pod template two containers",
			args: args{
				namespace:        "test",
				containerToDebug: "test",
				debugimage:       "test:latest",
				secretname:       "secret",
			},
			inputfile:  "testdata/add-pod-template/podtemplate_test_two_container.yaml",
			goldenfile: "testdata/add-pod-template/podtemplate_test_two_container.golden",
			wantErr:    false,
		},
		{
			name: "Add debug container to pod template where debug container already exists",
			args: args{
				namespace:        "test",
				containerToDebug: "test",
				debugimage:       "test:latest",
				secretname:       "secret",
			},
			inputfile: "testdata/add-pod-template/podtemplate_test_container_exists.yaml",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := unMarshalInputfile(tt.inputfile)
			actual, err := AddDebugContainerPodTemplate(input, tt.args.namespace, tt.args.containerToDebug, tt.args.debugimage, tt.args.secretname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error, no error returned")
				return
			}
			var expected corev1.PodTemplateSpec
			err = readUpdateGoldenFile(tt.goldenfile, *update, &expected, actual)
			eS := podTemplateToString(expected)
			aS := podTemplateToString(actual)
			if eS != aS {
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(podTemplateToString(expected), podTemplateToString(actual), false)
				t.Errorf("Returned PodTemplateSpec did not match. Difference:\n%s", dmp.DiffPrettyText(diffs))
			}
		})
	}
}

func TestRemoveDebugContainerPodTemplate(t *testing.T) {
	type args struct {
		namespace        string
		containerToDebug string
	}
	tests := []struct {
		name       string
		args       args
		inputfile  string
		goldenfile string
		wantErr    bool
	}{
		{
			name: "Remove debug container to pod template",
			args: args{
				namespace:        "test",
				containerToDebug: "",
			},
			inputfile:  "testdata/remove-pod-template/podtemplate_test_one_container.yaml",
			goldenfile: "testdata/remove-pod-template/podtemplate_test_one_container.golden",
			wantErr:    false,
		},
		{
			name: "Remove debug container to pod template two containers",
			args: args{
				namespace:        "test",
				containerToDebug: "test",
			},
			inputfile:  "testdata/remove-pod-template/podtemplate_test_two_container.yaml",
			goldenfile: "testdata/remove-pod-template/podtemplate_test_two_container.golden",
			wantErr:    false,
		},
		{
			name: "Remove debug container to pod template where debug container not exists",
			args: args{
				namespace:        "test",
				containerToDebug: "test",
			},
			inputfile: "testdata/remove-pod-template/podtemplate_test_container_exists.yaml",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := unMarshalInputfile(tt.inputfile)
			actual, err := RemoveDebugContainerPodTemplate(input, tt.args.namespace, tt.args.containerToDebug)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error, no error returned")
				return
			}
			var expected corev1.PodTemplateSpec
			err = readUpdateGoldenFile(tt.goldenfile, *update, &expected, actual)
			eS := podTemplateToString(expected)
			aS := podTemplateToString(actual)
			if eS != aS {
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(podTemplateToString(expected), podTemplateToString(actual), false)
				t.Errorf("Returned PodTemplateSpec did not match. Difference:\n%s", dmp.DiffPrettyText(diffs))
			}
		})
	}
}

func unMarshalInputfile(inputfile string) (corev1.PodTemplateSpec, error) {
	var template corev1.PodTemplateSpec
	f, err := os.Open(inputfile)
	if err != nil {
		return template, fmt.Errorf("failed to open inputfile %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return template, fmt.Errorf("failed to read inputfile %v", err)
	}
	err = yaml.Unmarshal(b, &template)
	return template, err
}

func readUpdateGoldenFile(goldenfile string, update bool, expected interface{}, actual interface{}) error {
	if update {
		data, err := yaml.Marshal(actual)
		if err != nil {
			return fmt.Errorf("failed to marhal input %v", err)
		}
		err = ioutil.WriteFile(goldenfile, data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write goldefile %v", err)
		}
	}
	f, err := os.Open(goldenfile)
	if err != nil {
		return fmt.Errorf("failed to open goldenfile %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read goldenfile %v", err)
	}
	err = yaml.Unmarshal(b, expected)
	return err
}

func podTemplateToString(template corev1.PodTemplateSpec) string {
	data, err := yaml.Marshal(template)
	if err != nil {
		return ""
	}
	return string(data)
}

func TestDDConfigFromPodTemplate(t *testing.T) {
	type args struct {
		template corev1.PodTemplateSpec
	}
	tests := []struct {
		name    string
		args    args
		want    DDConfig
		wantErr bool
	}{
		{
			name: "Test DDConfig from pod template",
			args: args{
				template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"dev.local/dd-added": "true",
							"dev.local/dd-apply": `{"containerToDebug":"test","debugContainerName":"debug","tmpdirAdded":false,"secretMount":"secret"}`,
						},
					},
				},
			},
			want: DDConfig{
				ContainerToDebug:   "test",
				DebugContainerName: "debug",
				TmpdirAdded:        false,
				SecretName:         "secret",
			},
			wantErr: false,
		},
		{
			name: "Test DDConfig config not applied",
			args: args{
				template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Test DDConfig config not valid",
			args: args{
				template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"dev.local/dd-added": "true",
							"dev.local/dd-apply": `{"containerToDebug":test,"debugContainerName":"debug","tmpdirAdded":false,"secretMount":"secret"}`,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DDConfigFromPodTemplate(tt.args.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("DDConfigFromPodTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DDConfigFromPodTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

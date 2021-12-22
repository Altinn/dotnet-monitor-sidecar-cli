package errors

import (
	"fmt"
	"testing"
)

func TestIsAlreadyPresent(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "Returns true if the error has value 'debug sidecar already present'",
			err:  fmt.Errorf("debug sidecar already present"),
			want: true,
		},
		{
			name: "Returns false if the error does not have value 'debug sidecar already present'",
			err:  fmt.Errorf("debug sidecar not present"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlreadyPresent(tt.err); got != tt.want {
				t.Errorf("IsAlreadyPresent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotPresent(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "Returns true if the error has value 'debug sidecar not present'",
			err:  fmt.Errorf("debug sidecar not present"),
			want: true,
		},
		{
			name: "Returns false if the error does not have value 'debug sidecar not present'",
			err:  fmt.Errorf("debug sidecar already present"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotPresent(tt.err); got != tt.want {
				t.Errorf("IsAlreadyPresent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "Returns true if the error has value 'resource not found'",
			err:  fmt.Errorf("resource not found"),
			want: true,
		},
		{
			name: "Returns false if the error does not have value 'resource not found'",
			err:  fmt.Errorf("resource found"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

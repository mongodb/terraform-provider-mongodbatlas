package cluster_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cluster"
)

func TestIsChangeStreamOptionsMinRequiredMajorVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty input", "", true},
		{"Valid input equal to 6", "6", true},
		{"Valid input greater than 6", "7.0", true},
		{"Valid input less than 6", "5", false},
		{"Valid float input greater", "6.5", true},
		{"Valid float input less", "5.9", false},
		{"Valid float complete semantic version", "6.0.2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cluster.IsChangeStreamOptionsMinRequiredMajorVersion(&tt.input); got != tt.want {
				t.Errorf("abc(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

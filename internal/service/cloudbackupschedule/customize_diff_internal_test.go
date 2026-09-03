package cloudbackupschedule

import "testing"

func TestShouldForceEmptyCopySettings(t *testing.T) {
	tests := []struct {
		name            string
		autoCopyEnabled bool
		rawEmpty        bool
		want            bool
	}{
		{
			name:            "auto false and raw empty forces empty",
			autoCopyEnabled: false,
			rawEmpty:        true,
			want:            true,
		},
		{
			name:            "auto true and raw empty skips force",
			autoCopyEnabled: true,
			rawEmpty:        true,
			want:            false,
		},
		{
			name:            "auto true and raw non-empty does not force",
			autoCopyEnabled: true,
			rawEmpty:        false,
			want:            false,
		},
		{
			name:            "auto false and raw non-empty does not force",
			autoCopyEnabled: false,
			rawEmpty:        false,
			want:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldForceEmptyCopySettings(tt.autoCopyEnabled, tt.rawEmpty)
			if got != tt.want {
				t.Errorf("shouldForceEmptyCopySettings(%v, %v) = %v, want %v", tt.autoCopyEnabled, tt.rawEmpty, got, tt.want)
			}
		})
	}
}

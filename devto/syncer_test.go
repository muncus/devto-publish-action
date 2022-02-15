// Syncer allows for updating of existing dev.to articles from a filesystem.

package devto

import (
	"testing"
)

func TestSyncer_LoadState(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]SyncerStateRecord
	}{
		{
			name:  "old format",
			input: "{\"foo\": 1}",
			expected: map[string]SyncerStateRecord{
				"foo": {Id: 1},
			},
		},
		{
			name:  "new format",
			input: `{"foo": { "id": 2}}`,
			expected: map[string]SyncerStateRecord{
				"foo": {Id: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewEmptySyncer()
			if err := s.LoadState([]byte(tt.input)); err != nil {
				t.Errorf("Syncer.LoadState() error = %v", err)
			}
			for k, v := range tt.expected {
				if s.StateMap[k] != v {
					t.Errorf("%s(%s): got %v, want %v", tt.name, k, s.StateMap[k], v)
				}
			}
		})
	}
}

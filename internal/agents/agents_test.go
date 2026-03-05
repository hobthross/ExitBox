package agents

import (
	"testing"
)

func TestRegistryGetAll(t *testing.T) {
	for _, name := range Names() {
		a := Get(name)
		if a == nil {
			t.Errorf("Get(%q) returned nil", name)
			continue
		}
		if a.Name() != name {
			t.Errorf("Get(%q).Name() = %q", name, a.Name())
		}
	}
}

func TestGetUnknown(t *testing.T) {
	if a := Get("nonexistent"); a != nil {
		t.Errorf("Get(\"nonexistent\") should return nil, got %v", a)
	}
}


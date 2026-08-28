package model

import "testing"

func TestLayoutValidateRejectsUnknownModule(t *testing.T) {
	layout := DefaultLayout()
	layout.Modules[1] = "metrics"

	if err := layout.Validate(); err == nil {
		t.Fatal("expected an unknown module to be rejected")
	}
}

func TestLayoutValidateAcceptsValidCustomization(t *testing.T) {
	layout := Layout{
		Modules: []string{"ai-lab", "overview", "resources", "posts"},
		Hidden:  []string{"posts"},
		Density: "compact",
		Theme:   "dark",
	}

	if err := layout.Validate(); err != nil {
		t.Fatalf("expected valid layout, got %v", err)
	}
}

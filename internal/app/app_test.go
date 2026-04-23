package app

import "testing"

func TestSwitchScreenForTest(t *testing.T) {
	model := New(BuildInfo{Version: "test", Commit: "none", Date: "unknown"})
	model = model.SwitchScreenForTest("settings")

	if model.CurrentScreenID() != "settings" {
		t.Fatalf("expected settings screen, got %q", model.CurrentScreenID())
	}
}

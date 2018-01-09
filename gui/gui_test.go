package gui


import (
	"testing"
)

func TestNewGui(t *testing.T) {
	gui := NewGui()

	if gui.ddcci == nil {
		t.Error("Expected Gui has ddcci initialized")
	}
}

func TestMonitorListing(t *testing.T) {

}
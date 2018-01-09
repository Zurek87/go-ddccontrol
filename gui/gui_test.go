package gui


import (
	"testing"
	"time"
)

func TestNewGui(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Got error:", r)
		}
	}()

	gui := NewGui()

	if gui.ddcci == nil {
		t.Error("Expected Gui has ddcci initialized")
	}
	gui.Show()
}

func TestChangeMonitor(t *testing.T) {
	gui := NewGui()
	gui.Show()

	list := gui.ddcci.MonitorList()

	// select first monitor:
	err := gui.SelectMonitor(0)
	if err != nil {
		t.Error("Expected Select monitor '0'; Got error:", err)
	}
	gui.selected.SetBrightness(100)
	count := len(list)
	if count > 1 {
		err = gui.SelectMonitor(count - 1)
		if err != nil {
			t.Errorf("Expected Select monitor '%v'; Got error: %v", count - 1, err)
		}
		gui.selected.SetBrightness(100)
		time.Sleep(1 * time.Second)
		gui.selected.SetBrightness(0)
		gui.SelectMonitor(0)
	} else {
		t.Skip("Only one monitor found")
	}
	gui.selected.SetBrightness(0)
}
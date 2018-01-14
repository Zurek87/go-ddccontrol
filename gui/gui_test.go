package gui


import (
	"testing"
	"time"
	"github.com/zurek87/go-gtk3/glib"
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
	defer func() {
		if r := recover(); r != nil {
			t.Error("Got error:", r)
		}
	}()
	gui := NewGui()
	gui.Show()

	list := gui.ddcci.MonitorList()

	// select first monitor:
	err := gui.SelectMonitor(0)
	if err != nil {
		t.Error("Expected Select monitor '0'; Got error:", err)
	}
	gui.SetBrightness(100)
	count := len(list)
	if count > 1 {
		err = gui.SelectMonitor(count - 1)
		if err != nil {
			t.Errorf("Expected Select monitor '%v'; Got error: %v", count - 1, err)
		}
		gui.SetBrightness(100)
		time.Sleep(1 * time.Second)
		gui.SetBrightness(0)
		gui.SelectMonitor(0)
	} else {
		t.Skip("Only one monitor found")
	}
	gui.SetBrightness(0)
}

func TestAddMenuItems(t *testing.T) {
	gui := NewGui()

	ddcciMonitorList := gui.ddcci.MonitorList()

	if len(ddcciMonitorList) == 0 {
		t.Error("Empty monitor list from DDCci")
	}
	list := guiMonitorList(ddcciMonitorList)
	if len(list) == 0 {
		t.Error("Empty monitor list from GUI")
	}
	// add select monitor
	group := glib.GSListAlloc()
	for _, info := range list{
		item , gr := createChoseMonitorMenuItem(info, group, &gui)
		if item == nil {
			t.Error("MenuItem not created")
		}
		if gr == nil {
			t.Error("No group returned")
		}
		group = gr // go-gtk bug
	}

	// add config monitor
	for _, info := range list{
		item := createConfigMonitorMenuItem(info)
		if item == nil {
			t.Error("MenuItem not created")
		}
		sub := item.GetSubmenu()
		if sub == nil {
			t.Error("MenuItem not have SubMenu")
		}
	}
}
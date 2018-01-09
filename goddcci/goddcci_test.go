package goddcci

import (
	"testing"
)

func TestInitDDCci2(t *testing.T) {
	err := initDDCci()
	if err != nil{
		t.Error("Expected Success","Got error:", err)
	}
}


func TestDDCci(t *testing.T) {
	err := initDDCci()
	if err != nil {
		t.Error("Expected initialized success","Got error:", err)
	}

	monitorList, err := probeDDCci()
	if err != nil {
		t.Error("Expected monitorList, got :", monitorList, "with error:", err)
	}
	ddcci := makeDDCci(monitorList)
	ddcci.detectSupportedMonitors()
	if ddcci.count == 0 {
		t.Error("Expected to find supported monitor")
	}
	err = ddcci.openMonitor(0)
	if err != nil {
		t.Error("Expected to open monitor, got error:", err)
	}

	defaultMonitor, err := ddcci.DefaultMonitor()
	if err != nil {
		t.Error("Expected get default monitor, got error:", err)
	}
	if defaultMonitor.Name == ""{
		t.Error("Expected to first monitor have name")
	}

	list := ddcci.MonitorList()
	if len(list) == 0 {
		t.Error("Empty monitor list")
	}
}
package goddcci

import (
	"testing"
	"time"
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

	err = defaultMonitor.SetBrightness(10)
	if err != nil {
		t.Error("Expected to set brightness, got error:", err)
	} else {
		time.Sleep(1 * time.Second)
		defaultMonitor.SetBrightness(0)
	}
}
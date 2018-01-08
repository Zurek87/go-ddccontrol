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
	if len(ddcci.supported) == 0 {
		t.Error("Expected to find supported monitor")
	}
	err = ddcci.openMonitor()
	if err != nil {
		t.Error("Expected to open monitor, got error:", err)
	}
}